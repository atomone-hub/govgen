package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	appparams "github.com/atomone-hub/govgen/v1/app/params"
	"github.com/atomone-hub/govgen/v1/x/gov/simulation"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func TestProposalContents(t *testing.T) {
	// initialize parameters
	s := rand.NewSource(1)
	r := rand.New(s)

	ctx := sdk.NewContext(nil, tmproto.Header{}, true, nil)
	accounts := simtypes.RandomAccounts(r, 3)

	// execute ProposalContents function
	weightedProposalContent := simulation.ProposalContents()
	require.Len(t, weightedProposalContent, 1)

	w0 := weightedProposalContent[0]

	// tests w0 interface:
	require.Equal(t, simulation.OpWeightMsgDeposit, w0.AppParamsKey())
	require.Equal(t, appparams.DefaultWeightTextProposal, w0.DefaultWeight())

	content := w0.ContentSimulatorFn()(r, ctx, accounts)

	require.Equal(t, "NxImpptHBIFDQfnxaTiOBJUgNzvqHbVcVJYlIFWFlzFqqRTTyFzDUMntPzyRamUFqeJAEaSHIuUHZoTWDjWXsYxYvwXwXZEsjRQKgKMselyUqWXMbHzRNDHnMzhWSirUgVggjiBxtWDfhzPDgrorEoNmDEiDdBldYegphCBTYWrmFFXNjxhtygsGBFHTejaKjMsqNdikEzDalEyWRHfJhKqifCKsedVuuJbQMbmRVuIPDluAWGpngjgBjOxuRFwSadayHNIhVVmNWBbfaTOldclxTTLUMvaBnLfwjHTtsKetEIvgrxLijhKJNablmvqpWIWsmhWQAYNLycREypoASHnyKWrxpoNLBJuyCGysZJgXbQAAmSIbGxMFXuwMVGZgBiZWfPWorAfjBeekCFvljHAtVZaTOsRxbPIioNxLTnWUTzGTvaNhplQQPmMADRRDuUIsiBpnGqPheKmLnopieVseFdTSAvOCacxaqFWFuXzsrVZzlGfeRpClwKuGEBujaPrzSLjVIOMvLlWxuznEOXlxbZroBRVEvEfBBAHOECribZNrYiFnzQqQmBnLksmFNAadusWAGltuqYNntgOlgOGwSdDjWdLboWyAWIcCfmpGJTfbljKPriLehwObuszICkaXNUkmeddeeRulbZBXJVLgteiKIfofGdNBregwUPlINQECatDSNXSIuefyMxxoKfcmjHEwbVtFiXtEnLJkLHUghmzFiymrgBChucZgOQUpGGVQEpRtIQjIBxYhtZPgUORdxXNWUMErWrUeriqYJPcgIDgLMWAyuuQnsHncCtjvHmvFbzYErxeunQllYDUVlXaRBveRUKeXwEGJFTSAqZtaBSDGDtzlADCnGjuTmYMJlapRsWfugmjwKEuoXJVpZvlcHeFvVvRRktRVGwzLfKezPEMABZtbLExQIjynSoahmkmoTHefdzFoBHMcQHFkKVHhpNtudPqJrYuQswzFuFHbSmpNltFnYJpvMrAYHFrNouZaanEUGHvbHIUUFTCtZrcpRHwgjblxlDNJWzHdBNpAXKJPHWQdrGYcAHSctgVlqwqHoLfHsXUdStwfefwzqLuKEhmMyYLdbZrcPgYqjNHxPexsruwEGStAneKbWkQDDIlCWBLSiAASNhZqNFlPtfqPJoxKsgMdzjWqLWdqKQuJqWPMvwPQWZUtVMOTMYKJbfdlZsjdsomuScvDmbDkgRualsxDvRJuCAmPOXitIbcyWsKGSdrEunFAOdmXnsuyFVgJqEjbklvmwrUlsxjRSfKZxGcpayDdgoFcnVSutxjRgOSFzPwidAjubMncNweqpbxhXGchpZUxuFDOtpnhNUycJICRYqsPhPSCjPTWZFLkstHWJxvdPEAyEIxXgLwbNOjrgzmaujiBABBIXvcXpLrbcEWNNQsbjvgJFgJkflpRohHUutvnaUqoopuKjTDaemDeSdqbnOzcfJpcTuAQtZoiLZOoAIlboFDAeGmSNwkvObPRvRWQgWkGkxwtPauYgdkmypLjbqhlHJIQTntgWjXwZdOyYEdQRRLfMSdnxqppqUofqLbLQDUjwKVKfZJUJQPsWIPwIVaSTrmKskoAhvmZyJgeRpkaTfGgrJzAigcxtfshmiDCFkuiluqtMOkidknnTBtumyJYlIsWLnCQclqdVmikUoMOPdPWwYbJxXyqUVicNxFxyqJTenNblyyKSdlCbiXxUiYUiMwXZASYfvMDPFgxniSjWaZTjHkqlJvtBsXqwPpyVxnJVGFWhfSxgOcduoxkiopJvFjMmFabrGYeVtTXLhxVUEiGwYUvndjFGzDVntUvibiyZhfMQdMhgsiuysLMiePBNXifRLMsSmXPkwlPloUbJveCvUlaalhZHuvdkCnkSHbMbmOnrfEGPwQiACiPlnihiaOdbjPqPiTXaHDoJXjSlZmltGqNHHNrcKdlFSCdmVOuvDcBLdSklyGJmcLTbSFtALdGlPkqqecJrpLCXNPWefoTJNgEJlyMEPneVaxxduAAEqQpHWZodWyRkDAxzyMnFMcjSVqeRXLqsNyNtQBbuRvunZflWSbbvXXdkyLikYqutQhLPONXbvhcQZJPSWnOulqQaXmbfFxAkqfYeseSHOQidHwbcsOaMnSrrmGjjRmEMQNuknupMxJiIeVjmgZvbmjPIQTEhQFULQLBMPrxcFPvBinaOPYWGvYGRKxLZdwamfRQQFngcdSlvwjfaPbURasIsGJVHtcEAxnIIrhSriiXLOlbEBLXFElXJFGxHJczRBIxAuPKtBisjKBwfzZFagdNmjdwIRvwzLkFKWRTDPxJCmpzHUcrPiiXXHnOIlqNVoGSXZewdnCRhuxeYGPVTfrNTQNOxZmxInOazUYNTNDgzsxlgiVEHPKMfbesvPHUqpNkUqbzeuzfdrsuLDpKHMUbBMKczKKWOdYoIXoPYtEjfOnlQLoGnbQUCuERdEFaptwnsHzTJDsuZkKtzMpFaZobynZdzNydEeJJHDYaQcwUxcqvwfWwNUsCiLvkZQiSfzAHftYgAmVsXgtmcYgTqJIawstRYJrZdSxlfRiqTufgEQVambeZZmaAyRQbcmdjVUZZCgqDrSeltJGXPMgZnGDZqISrGDOClxXCxMjmKqEPwKHoOfOeyGmqWqihqjINXLqnyTesZePQRqaWDQNqpLgNrAUKulklmckTijUltQKuWQDwpLmDyxLppPVMwsmBIpOwQttYFMjgJQZLYFPmxWFLIeZihkRNnkzoypBICIxgEuYsVWGIGRbbxqVasYnstWomJnHwmtOhAFSpttRYYzBmyEtZXiCthvKvWszTXDbiJbGXMcrYpKAgvUVFtdKUfvdMfhAryctklUCEdjetjuGNfJjajZtvzdYaqInKtFPPLYmRaXPdQzxdSQfmZDEVHlHGEGNSPRFJuIfKLLfUmnHxHnRjmzQPNlqrXgifUdzAGKVabYqvcDeYoTYgPsBUqehrBhmQUgTvDnsdpuhUoxskDdppTsYMcnDIPSwKIqhXDCIxOuXrywahvVavvHkPuaenjLmEbMgrkrQLHEAwrhHkPRNvonNQKqprqOFVZKAtpRSpvQUxMoXCMZLSSbnLEFsjVfANdQNQVwTmGxqVjVqRuxREAhuaDrFgEZpYKhwWPEKBevBfsOIcaZKyykQafzmGPLRAKDtTcJxJVgiiuUkmyMYuDUNEUhBEdoBLJnamtLmMJQgmLiUELIhLpiEvpOXOvXCPUeldLFqkKOwfacqIaRcnnZvERKRMCKUkMABbDHytQqQblrvoxOZkwzosQfDKGtIdfcXRJNqlBNwOCWoQBcEWyqrMlYZIAXYJmLfnjoJepgSFvrgajaBAIksoyeHqgqbGvpAstMIGmIhRYGGNPRIfOQKsGoKgxtsidhTaAePRCBFqZgPDWCIkqOJezGVkjfYUCZTlInbxBXwUAVRsxHTQtJFnnpmMvXDYCVlEmnZBKhmmxQOIQzxFWpJQkQoSAYzTEiDWEOsVLNrbfzeHFRyeYATakQQWmFDLPbVMCJcWjFGJjfqCoVzlbNNEsqxdSmNPjTjHYOkuEMFLkXYGaoJlraLqayMeCsTjWNRDPBywBJLAPVkGQqTwApVVwYAetlwSbzsdHWsTwSIcctkyKDuRWYDQikRqsKTMJchrliONJeaZIzwPQrNbTwxsGdwuduvibtYndRwpdsvyCktRHFalvUuEKMqXbItfGcNGWsGzubdPMYayOUOINjpcFBeESdwpdlTYmrPsLsVDhpTzoMegKrytNVZkfJRPuDCUXxSlSthOohmsuxmIZUedzxKmowKOdXTMcEtdpHaPWgIsIjrViKrQOCONlSuazmLuCUjLltOGXeNgJKedTVrrVCpWYWHyVrdXpKgNaMJVjbXxnVMSChdWKuZdqpisvrkBJPoURDYxWOtpjzZoOpWzyUuYNhCzRoHsMjmmWDcXzQiHIyjwdhPNwiPqFxeUfMVFQGImhykFgMIlQEoZCaRoqSBXTSWAeDumdbsOGtATwEdZlLfoBKiTvodQBGOEcuATWXfiinSjPmJKcWgQrTVYVrwlyMWhxqNbCMpIQNoSMGTiWfPTCezUjYcdWppnsYJihLQCqbNLRGgqrwHuIvsazapTpoPZIyZyeeSueJuTIhpHMEJfJpScshJubJGfkusuVBgfTWQoywSSliQQSfbvaHKiLnyjdSbpMkdBgXepoSsHnCQaYuHQqZsoEOmJCiuQUpJkmfyfbIShzlZpHFmLCsbknEAkKXKfRTRnuwdBeuOGgFbJLbDksHVapaRayWzwoYBEpmrlAxrUxYMUekKbpjPNfjUCjhbdMAnJmYQVZBQZkFVweHDAlaqJjRqoQPoOMLhyvYCzqEuQsAFoxWrzRnTVjStPadhsESlERnKhpEPsfDxNvxqcOyIulaCkmPdambLHvGhTZzysvqFauEgkFRItPfvisehFmoBhQqmkfbHVsgfHXDPJVyhwPllQpuYLRYvGodxKjkarnSNgsXoKEMlaSKxKdcVgvOkuLcfLFfdtXGTclqfPOfeoVLbqcjcXCUEBgAGplrkgsmIEhWRZLlGPGCwKWRaCKMkBHTAcypUrYjWwCLtOPVygMwMANGoQwFnCqFrUGMCRZUGJKTZIGPyldsifauoMnJPLTcDHmilcmahlqOELaAUYDBuzsVywnDQfwRLGIWozYaOAilMBcObErwgTDNGWnwQMUgFFSKtPDMEoEQCTKVREqrXZSGLqwTMcxHfWotDllNkIJPMbXzjDVjPOOjCFuIvTyhXKLyhUScOXvYthRXpPfKwMhptXaxIxgqBoUqzrWbaoLTVpQoottZyPFfNOoMioXHRuFwMRYUiKvcWPkrayyTLOCFJlAyslDameIuqVAuxErqFPEWIScKpBORIuZqoXlZuTvAjEdlEWDODFRregDTqGNoFBIHxvimmIZwLfFyKUfEWAnNBdtdzDmTPXtpHRGdIbuucfTjOygZsTxPjf", content.GetDescription())
	require.Equal(t, "XhSUkMhPjMaxKlMIJMOXcnQfyzeOcbWwNbeHVIkPZBSpYuLyYggwexjxusrBqDOTtGTOWeLrQKjLxzIivHSlcxgdXhhuTSkuxKGLwQvuyNhYFmBZHeAerqyNEUzXPFGkqEGqiQWIXnku", content.GetTitle())
	require.Equal(t, "gov", content.ProposalRoute())
	require.Equal(t, "Text", content.ProposalType())
}
