package api

import (
	"encoding/json"
	"github.com/ixre/gof/api"
	"go2o/core/infrastructure/domain"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/thrift"
	"strconv"
	"strings"
)

var _ api.Handler = new(MemberApi)

var provider = map[string]string{
	"alipay": "支付宝",
	"wepay":  "微信支付",
	"unipay": "云闪付",
}

type MemberApi struct {
	*apiUtil
}

func (m MemberApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"login":           m.login,
		"get":             m.getMember,
		"account":         m.account,
		"profile":         m.profile,
		"checkToken":      m.checkToken,
		"complex":         m.complex,
		"bankcard":        m.bankcard,
		"receipts_code":   m.receiptsCode,
		"save_receipts":   m.saveReceiptsCode,
		"toggle_receipts": m.toggleReceipts,
	})
}

// 登录
func (m MemberApi) login(ctx api.Context) interface{} {
	form := ctx.Form()
	user := strings.TrimSpace(form.GetString("user"))
	pwd := strings.TrimSpace(form.GetString("pwd"))
	if len(user) == 0 || len(pwd) == 0 {
		return api.ResponseWithCode(2, "缺少参数: user or pwd")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err != nil {
		return api.ResponseWithCode(3, "网络连接失败")
	}
	defer trans.Close()
	encPwd := domain.MemberSha1Pwd(pwd)
	r, _ := cli.CheckLogin(thrift.Context, user, encPwd, true)
	if r.ErrCode == 0 {
		memberId, _ := strconv.Atoi(r.Data["id"])
		token, _ := cli.GetToken(thrift.Context, int64(memberId), true)
		r.Data["token"] = token
		return r
	} else {
		return api.ResponseWithCode(int(r.ErrCode), r.ErrMsg)
	}
}

// 注册
func (m MemberApi) Register(ctx api.Context) interface{} {
	user := ctx.Form().GetString("user")
	pwd := ctx.Form().GetString("pwd")
	phone := ctx.Form().GetString("phone")
	regFrom := ctx.Form().GetString("reg_from")       // 注册来源
	inviteCode := ctx.Form().GetString("invite_code") // 邀请码
	regIp := ctx.Form().GetString("$user_ip_addr")    // IP地址
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mp := map[string]string{
			"reg_ip":      regIp,
			"reg_from":    regFrom,
			"invite_code": inviteCode,
		}
		r, _ := cli.RegisterMemberV2(thrift.Context, user, pwd, 0, "", phone, "", "", mp)
		return r
	}
	return m.SResult(err)
}

// 账号信息
func (m MemberApi) account(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, err1 := cli.GetAccount(thrift.Context, int64(memberId))
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 账号信息
func (m MemberApi) complex(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, _ := cli.Complex(thrift.Context, memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}

// 银行卡
func (m MemberApi) bankcard(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, _ := cli.Bankcards(thrift.Context, memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}

// 账号信息
func (m MemberApi) profile(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, err1 := cli.GetProfile(thrift.Context, memberId)
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 账号信息
func (m MemberApi) checkToken(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	token := strings.TrimSpace(ctx.Form().GetString("token"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code or token")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		r, err1 := cli.CheckToken(thrift.Context, memberId, token)
		if err1 == nil {
			return r
		}
		err = err1
	}
	return api.NewErrorResponse(err.Error())
}

// 获取会员信息
func (m MemberApi) getMember(ctx api.Context) interface{} {
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	if len(code) == 0 {
		return api.NewErrorResponse("missing params: code")
	}
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
		if memberId <= 0 {
			return api.NewErrorResponse("no such member")
		}
		r, _ := cli.GetMember(thrift.Context, memberId)
		return r
	}
	return api.NewErrorResponse(err.Error())
}

func (m MemberApi) receiptsCode(ctx api.Context) interface{} {
	trans, cli, _ := thrift.MemberServeClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	arr, _ := cli.ReceiptsCodes(thrift.Context, memberId)
	mp := map[string]interface{}{
		"list":     arr,
		"provider": provider,
	}
	return mp
}
func (m MemberApi) saveReceiptsCode(ctx api.Context) interface{} {
	trans, cli, _ := thrift.MemberServeClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	data := ctx.Form().GetBytes("data")
	c := &member_service.SReceiptsCode{}
	json.Unmarshal(data, c)
	if _, ok := provider[c.Identity]; !ok {
		return api.NewErrorResponse("不支持的收款码")
	}
	if c.ID == 0 {
		c.State = 1
	}
	memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	r, _ := cli.SaveReceiptsCode(thrift.Context, memberId, c)
	return r
}
func (m MemberApi) toggleReceipts(ctx api.Context) interface{} {
	trans, cli, _ := thrift.MemberServeClient()
	defer trans.Close()
	code := strings.TrimSpace(ctx.Form().GetString("code"))
	id := ctx.Form().GetInt("id")
	memberId, _ := cli.SwapMemberId(thrift.Context, member_service.ECredentials_Code, code)
	arr, _ := cli.ReceiptsCodes(thrift.Context, memberId)
	for _, v := range arr {
		if int(v.ID) == id {
			v.State = 1 - v.State
			r, _ := cli.SaveReceiptsCode(thrift.Context, memberId, v)
			return r
		}
	}
	return api.NewErrorResponse("no such data")
}
