package parser

import (
	"embed"
	"testing"

	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/internal/metrics"
	"github.com/altinity/transfer/pkg/abstract"
	parsersfactory "github.com/altinity/transfer/pkg/parsers"
	"github.com/altinity/transfer/pkg/stats"
	"github.com/altinity/transfer/tests/canon/parser/testcase"
	"github.com/altinity/transfer/tests/canon/validator"
	"github.com/stretchr/testify/require"
)

//go:embed samples/static/generic/*
var TestGenericSamples embed.FS

func TestGenericParsers(t *testing.T) {
	cases := testcase.LoadStaticTestCases(t, TestGenericSamples)

	for tc := range cases {
		t.Run(tc, func(t *testing.T) {
			currCase := cases[tc]
			parser, err := parsersfactory.NewParserFromParserConfig(currCase.ParserConfig, false, logger.Log, stats.NewSourceStats(metrics.NewRegistry().WithTags(map[string]string{
				"id": "TestParser_Do",
			})))
			require.NoError(t, err)
			require.NotNil(t, parser)
			res := parser.Do(currCase.Data, abstract.Partition{Topic: currCase.TopicName})
			require.NotNil(t, res)
			sink := validator.New(
				false,
				validator.ValuesTypeChecker,
				validator.Canonizator(t),
			)()
			require.NoError(t, sink.Push(res))
			require.NoError(t, sink.Close())
		})
	}
}
