package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/figment-networks/indexer-manager/structs"
	cStruct "github.com/figment-networks/indexer-manager/worker/connectivity/structs"
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/kava-worker/api/mapper"
	"github.com/figment-networks/kava-worker/api/types"
	"github.com/figment-networks/kava-worker/api/util"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var (
	errUnknownMessageType = fmt.Errorf("unknown message type")
)

// TxLogError Error message
type TxLogError struct {
	Codespace string  `json:"codespace"`
	Code      float64 `json:"code"`
	Message   string  `json:"message"`
}

// SearchTx is making search api call
func (c *Client) SearchTx(ctx context.Context, r structs.HeightHash, block structs.Block, perPage uint64) (txs []structs.Transaction, err error) {
	defer c.logger.Sync()

	numberOfItemsInBlock.Add(float64(block.NumberOfTransactions))
	page := uint64(1)
	for {
		now := time.Now()
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return txs, err
		}

		sCtx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		req, err := http.NewRequestWithContext(sCtx, http.MethodGet, c.baseURL+"/tx_search", nil)
		if err != nil {
			return txs, err
		}

		req.Header.Add("Content-Type", "application/json")
		if c.key != "" {
			req.Header.Add("Authorization", c.key)
		}

		q := req.URL.Query()
		s := strings.Builder{}
		s.WriteString(`"`)
		s.WriteString("tx.height=")
		s.WriteString(strconv.FormatUint(r.Height, 10))
		s.WriteString(`"`)

		q.Add("query", s.String())
		q.Add("page", strconv.FormatUint(page, 10))
		q.Add("per_page", strconv.FormatUint(perPage, 10))
		req.URL.RawQuery = q.Encode()

		resp, err := c.httpClient.Do(req)
		log.Debug("[COSMOS-API] Request Time (/tx_search)", zap.Duration("duration", time.Now().Sub(now)))
		if err != nil {
			return txs, err
		}

		if resp.StatusCode > 399 { // ERROR
			serverError, _ := ioutil.ReadAll(resp.Body)
			c.logger.Error("[COSMOS-API] error getting response from server", zap.Int("code", resp.StatusCode), zap.Any("response", string(serverError)))
			return txs, fmt.Errorf("error getting response from server %d %s", resp.StatusCode, string(serverError))
		}

		rawRequestHTTPDuration.WithLabels("/tx_search", resp.Status).Observe(time.Since(now).Seconds())

		decoder := json.NewDecoder(resp.Body)

		result := &types.GetTxSearchResponse{}
		if err = decoder.Decode(result); err != nil {
			c.logger.Error("[COSMOS-API] unable to decode result body", zap.Error(err))
			return txs, fmt.Errorf("unable to decode result body %w", err)
		}

		if result.Error.Message != "" {
			c.logger.Error("[COSMOS-API] Error getting search", zap.Any("result", result.Error.Message))
			return txs, fmt.Errorf("Error getting search: %s", result.Error.Message)
		}

		totalCount, err := strconv.ParseInt(result.Result.TotalCount, 10, 64)
		if err != nil {
			c.logger.Error("[COSMOS-API] Error getting totalCount", zap.Error(err), zap.Any("result", result), zap.String("query", req.URL.RawQuery), zap.Any("request", r))
			return txs, err
		}

		numberOfItemsInBlock.Add(float64(totalCount))
		c.logger.Debug("[COSMOS-API] Converting requests ", zap.Int("number", len(result.Result.Txs)))

		for _, txRaw := range result.Result.Txs {
			tx, err := rawToTransaction(ctx, txRaw, c.logger, c.cdc)
			if err != nil {
				return nil, err
			}
			//conversionDuration.WithLabels(txRaw.Tx.TypeUrl).Observe(time.Since(n).Seconds())
			tx.BlockHash = block.Hash
			tx.ChainID = block.ChainID
			tx.Time = block.Time
			txs = append(txs, tx)
		}

		if totalCount <= int64(len(txs)) {
			break
		}
		page++
	}

	c.logger.Debug("[COSMOS-API] Converted all requests ", zap.Int("number", len(txs)), zap.Uint64("height", r.Height))
	return txs, nil
}

// transform raw data from chain into transaction format with augmentation from blocks
func rawToTransaction(ctx context.Context, in types.TxResponse, logger *zap.Logger, cdc *codec.Codec) (trans structs.Transaction, err error) {
	defer logger.Sync()
	timer := metrics.NewTimer(transactionConversionDuration)
	tx := &auth.StdTx{}
	lf := []types.LogFormat{}
	txErrs := []TxLogError{}

	if err := json.Unmarshal([]byte(in.TxResult.Log), &lf); err != nil {
		if in.TxResult.Log != "" && in.TxResult.Code > 0 {
			txErrs = append(txErrs, TxLogError{
				Message:   in.TxResult.Log,
				Code:      in.TxResult.Code,
				Codespace: in.TxResult.Codespace,
			})
		}
	}

	for _, logf := range lf {
		tle := TxLogError{}
		if errin := json.Unmarshal([]byte(logf.Log), &tle); errin == nil && tle.Message != "" {
			txErrs = append(txErrs, tle)
		}
		s.WriteString(`"`)

		q.Add("query", s.String())
		q.Add("page", strconv.Itoa(page))
		q.Add("per_page", strconv.Itoa(perPage))
		req.URL.RawQuery = q.Encode()

		now := time.Now()
		resp, err := c.httpClient.Do(req)

	}

	txReader := strings.NewReader(in.TxData)
	base64Dec := base64.NewDecoder(base64.StdEncoding, txReader)
	_, err = cdc.UnmarshalBinaryLengthPrefixedReader(base64Dec, tx, 0)
	if err != nil {
		logger.Error("[COSMOS-API] Problem decoding raw transaction (cdc)", zap.Error(err), zap.Any("height", in.Height), zap.Any("raw_tx", in))
	}

	outTX := cStruct.OutResp{Type: "Transaction"}
	trans = structs.Transaction{
		Hash:   in.Hash,
		Memo:   tx.GetMemo(),
		RawLog: []byte(in.TxResult.Log),
	}

	for _, coin := range tx.Fee.Amount {
		trans.Fee = append(trans.Fee, structs.TransactionAmount{
			Text:     coin.Amount.String(),
			Numeric:  coin.Amount.BigInt(),
			Currency: coin.Denom,
		})
	}

	trans.Height, err = strconv.ParseUint(in.Height, 10, 64)
	if err != nil {
		outTX.Error = err
	}
	trans.GasWanted, err = strconv.ParseUint(in.TxResult.GasWanted, 10, 64)
	if err != nil {
		outTX.Error = err
	}
	trans.GasUsed, err = strconv.ParseUint(in.TxResult.GasUsed, 10, 64)
	if err != nil {
		outTX.Error = err
	}

	txReader.Seek(0, 0)
	trans.Raw = make([]byte, txReader.Len())
	txReader.Read(trans.Raw)

	presentIndexes := map[string]bool{}

	for index, msg := range tx.Msgs {
		tev := structs.TransactionEvent{
			ID: strconv.Itoa(index),
		}

		var ev structs.SubsetEvent
		var err error
		logAtIndex := findLog(lf, index)

		switch msg.Route() {
		case "auction":
			switch msg.Type() {
			case "place_bid":
				ev, err = mapper.AuctionPlaceBidToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "bank":
			switch msg.Type() {
			case "multisend":
				ev, err = mapper.BankMultisendToSub(msg, logAtIndex)
			case "send":
				ev, err = mapper.BankSendToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "bep3":
			switch msg.Type() {
			case "createAtomicSwap":
				ev, err = mapper.Bep3CreateAtomicSwapToSub(msg, logAtIndex)
			case "claimAtomicSwap":
				ev, err = mapper.Bep3ClaimAtomicSwapToSub(msg)
			case "refundAtomicSwap":
				ev, err = mapper.Bep3RefundAtomicSwapToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "cdp":
			switch msg.Type() {
			case "create_cdp":
				ev, err = mapper.CDPCreateCDPToSub(msg)
			case "deposit_cdp":
				ev, err = mapper.CDPDepositCDPToSub(msg)
			case "withdraw_cdp":
				ev, err = mapper.CDPWithdrawCDPToSub(msg)
			case "draw_cdp":
				ev, err = mapper.CDPDrawCDPToSub(msg)
			case "repay_cdp":
				ev, err = mapper.CDPRepayCDPToSub(msg)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "committee":
			switch msg.Type() {
			case "commmittee_submit_proposal":
				ev, err = mapper.CommitteeSubmitProposalToSub(msg)
			case "committee_vote":
				ev, err = mapper.CommitteeVoteToSub(msg)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "crisis":
			switch msg.Type() {
			case "verify_invariant":
				ev, err = mapper.CrisisVerifyInvariantToSub(msg)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "distribution":
			switch msg.Type() {
			case "withdraw_validator_commission":
				ev, err = mapper.DistributionWithdrawValidatorCommissionToSub(msg, logAtIndex)
			case "set_withdraw_address":
				ev, err = mapper.DistributionSetWithdrawAddressToSub(msg)
			case "withdraw_delegator_reward":
				ev, err = mapper.DistributionWithdrawDelegatorRewardToSub(msg, logAtIndex)
			case "fund_community_pool":
				ev, err = mapper.DistributionFundCommunityPoolToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "evidence":
			switch msg.Type() {
			case "submit_evidence":
				ev, err = mapper.EvidenceSubmitEvidenceToSub(msg)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "gov":
			switch msg.Type() {
			case "deposit":
				ev, err = mapper.GovDepositToSub(msg, logAtIndex)
			case "vote":
				ev, err = mapper.GovVoteToSub(msg)
			case "submit_proposal":
				ev, err = mapper.GovSubmitProposalToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "harvest":
			switch msg.Type() {
			case "harvest_deposit":
				ev, err = mapper.HarvestDepositToSub(msg, logAtIndex)
			case "harvest_withdraw":
				ev, err = mapper.HarvestWithdrawToSub(msg, logAtIndex)
			case "claim_harvest_reward":
				ev, err = mapper.HarvestClaimRewardToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "incentive":
			switch msg.Type() {
			case "claim_reward":
				ev, err = mapper.IncentiveClaimRewardToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "issuance":
			switch msg.Type() {
			case "issue_tokens":
				ev, err = mapper.IssuanceIssueTokensToSub(msg)
			case "redeem_tokens":
				ev, err = mapper.IssuanceRedeemTokensToSub(msg)
			case "block_address":
				ev, err = mapper.IssuanceBlockAddressToSub(msg)
			case "unblock_address":
				ev, err = mapper.IssuanceUnblockAddressToSub(msg)
			case "change_pause_status":
				ev, err = mapper.IssuanceMsgSetPauseStatusToSub(msg)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "pricefeed":
			switch msg.Type() {
			case "post_price":
				ev, err = mapper.PricefeedPostPrice(msg)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "slashing":
			switch msg.Type() {
			case "unjail":
				ev, err = mapper.SlashingUnjailToSub(msg)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		case "staking":
			switch msg.Type() {
			case "begin_unbonding":
				ev, err = mapper.StakingUndelegateToSub(msg, logAtIndex)
			case "edit_validator":
				ev, err = mapper.StakingEditValidatorToSub(msg)
			case "create_validator":
				ev, err = mapper.StakingCreateValidatorToSub(msg)
			case "delegate":
				ev, err = mapper.StakingDelegateToSub(msg, logAtIndex)
			case "begin_redelegate":
				ev, err = mapper.StakingBeginRedelegateToSub(msg, logAtIndex)
			default:
				err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
			}
		default:
			err = fmt.Errorf("problem with %s - %s: %w", msg.Route(), msg.Type(), errUnknownMessageType)
		}

		if len(ev.Type) > 0 {
			tev.Kind = msg.Type()
			tev.Sub = append(tev.Sub, ev)
		}

		if err != nil {
			if errors.Is(err, errUnknownMessageType) {
				unknownTransactions.WithLabels(msg.Route() + "/" + msg.Type()).Inc()
			} else {
				brokenTransactions.WithLabels(msg.Route() + "/" + msg.Type()).Inc()
			}

			logger.Error("[COSMOS-API] Problem decoding transaction ", zap.Error(err), zap.String("type", msg.Type()), zap.String("route", msg.Route()), zap.String("height", in.Height))
		}

		presentIndexes[tev.ID] = true
		trans.Events = append(trans.Events, tev)
	}

	for _, logf := range lf {
		msgIndex := strconv.FormatFloat(logf.MsgIndex, 'f', -1, 64)
		_, ok := presentIndexes[msgIndex]
		if ok {
			continue
		}

		tev := structs.TransactionEvent{
			ID: msgIndex,
		}
		for _, ev := range logf.Events {
			sub := structs.SubsetEvent{
				Type: []string{ev.Type},
			}
			for atk, attr := range ev.Attributes {
				sub.Module = attr.Module
				sub.Action = attr.Action

				if len(attr.Sender) > 0 {
					for _, senderID := range attr.Sender {
						sub.Sender = append(sub.Sender, structs.EventTransfer{Account: structs.Account{ID: senderID}})
					}
				}
				if len(attr.Recipient) > 0 {
					for _, recipientID := range attr.Recipient {
						sub.Recipient = append(sub.Recipient, structs.EventTransfer{Account: structs.Account{ID: recipientID}})
					}
				}
				if attr.CompletionTime != "" {
					cTime, _ := time.Parse(time.RFC3339Nano, attr.CompletionTime)
					sub.Completion = &cTime
				}
				if len(attr.Validator) > 0 {
					if sub.Node == nil {
						sub.Node = make(map[string][]structs.Account)
					}
					for k, v := range attr.Validator {
						w, ok := sub.Node[k]
						if !ok {
							w = []structs.Account{}
						}

						for _, validatorID := range v {
							w = append(w, structs.Account{ID: validatorID})
						}
						sub.Node[k] = w
					}
				}

				for index, amount := range attr.Amount {
					sliced := util.GetCurrency(amount)

					am := structs.TransactionAmount{
						Text: amount,
					}

					var (
						c       *big.Int
						exp     int32
						coinErr error
					)

					if len(sliced) == 3 {
						am.Currency = sliced[2]
						c, exp, coinErr = util.GetCoin(sliced[1])
					} else {
						c, exp, coinErr = util.GetCoin(amount)
					}

					if coinErr != nil {
						am.Numeric.Set(c)
						am.Exp = exp
					}

					if sub.Amount == nil {
						sub.Amount = make(map[string]structs.TransactionAmount)
					}
					sub.Amount[strconv.Itoa(index)] = am
				}
				ev.Attributes[atk] = nil
			}
			tev.Sub = append(tev.Sub, sub)
		}
		logf.Events = nil
		trans.Events = append(trans.Events, tev)
	}

	for _, txErr := range txErrs {
		if txErr.Message != "" {
			trans.Events = append(trans.Events, structs.TransactionEvent{
				Kind: "error",
				Sub: []structs.SubsetEvent{{
					Type:   []string{"error"},
					Module: txErr.Codespace,
					Error:  &structs.SubsetEventError{Message: txErr.Message},
				}},
			})
		}
	}

	numberOfItemsTransactions.Inc()
	outTX.Payload = trans
	timer.ObserveDuration()

	// GC Help
	lf = nil

	return trans, nil
}

// GetFromRaw returns raw data for plugin use;
func (c *Client) GetFromRaw(logger *zap.Logger, txReader io.Reader) []map[string]interface{} {
	tx := &auth.StdTx{}
	base64Dec := base64.NewDecoder(base64.StdEncoding, txReader)
	_, err := c.cdc.UnmarshalBinaryLengthPrefixedReader(base64Dec, tx, 0)
	if err != nil {
		logger.Error("[KAVA-API] Problem decoding raw transaction (cdc) ", zap.Error(err))
	}
	slice := []map[string]interface{}{}
	for _, coin := range tx.Fee.Amount {
		slice = append(slice, map[string]interface{}{
			"text":     coin.Amount.String(),
			"numeric":  coin.Amount.BigInt(),
			"currency": coin.Denom,
		})
	}
	return slice
}

func findLog(lf []types.LogFormat, index int) types.LogFormat {
	if len(lf) <= index {
		return types.LogFormat{}
	}
	if l := lf[index]; l.MsgIndex == float64(index) {
		return l
	}
	for _, l := range lf {
		if l.MsgIndex == float64(index) {
			return l
		}
	}
	return types.LogFormat{}
}
