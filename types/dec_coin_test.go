package types_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type decCoinTestSuite struct {
	suite.Suite
}

func TestDecCoinTestSuite(t *testing.T) {
	suite.Run(t, new(decCoinTestSuite))
}

func (s *decCoinTestSuite) TestNewDecCoin() {
	s.Require().NotPanics(func() {
		sdk.NewInt64DecCoin(testDenom1, 5)
	})
	s.Require().NotPanics(func() {
		sdk.NewInt64DecCoin(testDenom1, 0)
	})
	s.Require().NotPanics(func() {
		sdk.NewInt64DecCoin(strings.ToUpper(testDenom1), 5)
	})
	s.Require().Panics(func() {
		sdk.NewInt64DecCoin(testDenom1, -5)
	})
}

func (s *decCoinTestSuite) TestNewDecCoinFromDec() {
	s.Require().NotPanics(func() {
		sdk.NewDecCoinFromDec(testDenom1, sdk.NewDec(5))
	})
	s.Require().NotPanics(func() {
		sdk.NewDecCoinFromDec(testDenom1, sdk.ZeroDec())
	})
	s.Require().NotPanics(func() {
		sdk.NewDecCoinFromDec(strings.ToUpper(testDenom1), sdk.NewDec(5))
	})
	s.Require().Panics(func() {
		sdk.NewDecCoinFromDec(testDenom1, sdk.NewDec(-5))
	})
}

func (s *decCoinTestSuite) TestNewDecCoinFromCoin() {
	s.Require().NotPanics(func() {
		sdk.NewDecCoinFromCoin(sdk.Coin{testDenom1, sdk.NewInt(5)})
	})
	s.Require().NotPanics(func() {
		sdk.NewDecCoinFromCoin(sdk.Coin{testDenom1, sdk.NewInt(0)})
	})
	s.Require().NotPanics(func() {
		sdk.NewDecCoinFromCoin(sdk.Coin{strings.ToUpper(testDenom1), sdk.NewInt(5)})
	})
	s.Require().Panics(func() {
		sdk.NewDecCoinFromCoin(sdk.Coin{testDenom1, sdk.NewInt(-5)})
	})
}

func (s *decCoinTestSuite) TestDecCoinIsPositive() {
	dc := sdk.NewInt64DecCoin(testDenom1, 5)
	s.Require().True(dc.IsPositive())

	dc = sdk.NewInt64DecCoin(testDenom1, 0)
	s.Require().False(dc.IsPositive())
}

func (s *decCoinTestSuite) TestAddDecCoin() {
	decCoinA1 := sdk.NewDecCoinFromDec(testDenom1, sdk.NewDecWithPrec(11, 1))
	decCoinA2 := sdk.NewDecCoinFromDec(testDenom1, sdk.NewDecWithPrec(22, 1))
	decCoinB1 := sdk.NewDecCoinFromDec(testDenom2, sdk.NewDecWithPrec(11, 1))

	// regular add
	res := decCoinA1.Add(decCoinA1)
	s.Require().Equal(decCoinA2, res, "sum of coins is incorrect")

	// bad denom add
	s.Require().Panics(func() {
		decCoinA1.Add(decCoinB1)
	}, "expected panic on sum of different denoms")
}

func (s *decCoinTestSuite) TestAddDecCoins() {
	one := sdk.NewDec(1)
	zero := sdk.NewDec(0)
	two := sdk.NewDec(2)

	cases := []struct {
		inputOne sdk.DecCoins
		inputTwo sdk.DecCoins
		expected sdk.DecCoins
	}{
		{sdk.DecCoins{{testDenom1, one}, {testDenom2, one}}, sdk.DecCoins{{testDenom1, one}, {testDenom2, one}}, sdk.DecCoins{{testDenom1, two}, {testDenom2, two}}},
		{sdk.DecCoins{{testDenom1, zero}, {testDenom2, one}}, sdk.DecCoins{{testDenom1, zero}, {testDenom2, zero}}, sdk.DecCoins{{testDenom2, one}}},
		{sdk.DecCoins{{testDenom1, zero}, {testDenom2, zero}}, sdk.DecCoins{{testDenom1, zero}, {testDenom2, zero}}, sdk.DecCoins(nil)},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.Add(tc.inputTwo...)
		s.Require().Equal(tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
	}
}

func (s *decCoinTestSuite) TestFilteredZeroDecCoins() {
	cases := []struct {
		name     string
		input    sdk.DecCoins
		original string
		expected string
		panic    bool
	}{
		{
			name: "all greater than zero",
			input: sdk.DecCoins{
				{"testa", sdk.NewDec(1)},
				{"testb", sdk.NewDec(2)},
				{"testc", sdk.NewDec(3)},
				{"testd", sdk.NewDec(4)},
				{"teste", sdk.NewDec(5)},
			},
			original: "1.000000000000000000testa,2.000000000000000000testb,3.000000000000000000testc,4.000000000000000000testd,5.000000000000000000teste",
			expected: "1.000000000000000000testa,2.000000000000000000testb,3.000000000000000000testc,4.000000000000000000testd,5.000000000000000000teste",
			panic:    false,
		},
		{
			name: "zero coin in middle",
			input: sdk.DecCoins{
				{"testa", sdk.NewDec(1)},
				{"testb", sdk.NewDec(2)},
				{"testc", sdk.NewDec(0)},
				{"testd", sdk.NewDec(4)},
				{"teste", sdk.NewDec(5)},
			},
			original: "1.000000000000000000testa,2.000000000000000000testb,0.000000000000000000testc,4.000000000000000000testd,5.000000000000000000teste",
			expected: "1.000000000000000000testa,2.000000000000000000testb,4.000000000000000000testd,5.000000000000000000teste",
			panic:    false,
		},
		{
			name: "zero coin end (unordered)",
			input: sdk.DecCoins{
				{"teste", sdk.NewDec(5)},
				{"testc", sdk.NewDec(3)},
				{"testa", sdk.NewDec(1)},
				{"testd", sdk.NewDec(4)},
				{"testb", sdk.NewDec(0)},
			},
			original: "5.000000000000000000teste,3.000000000000000000testc,1.000000000000000000testa,4.000000000000000000testd,0.000000000000000000testb",
			expected: "1.000000000000000000testa,3.000000000000000000testc,4.000000000000000000testd,5.000000000000000000teste",
			panic:    false,
		},

		{
			name: "panic when same denoms in multiple coins",
			input: sdk.DecCoins{
				{"testa", sdk.NewDec(5)},
				{"testa", sdk.NewDec(3)},
				{"testa", sdk.NewDec(1)},
				{"testd", sdk.NewDec(4)},
				{"testb", sdk.NewDec(2)},
			},
			original: "5.000000000000000000teste,3.000000000000000000testc,1.000000000000000000testa,4.000000000000000000testd,0.000000000000000000testb",
			expected: "1.000000000000000000testa,3.000000000000000000testc,4.000000000000000000testd,5.000000000000000000teste",
			panic:    true,
		},
	}

	for _, tt := range cases {
		if tt.panic {
			s.Require().Panics(func() { sdk.NewDecCoins(tt.input...) }, "Should panic due to multiple coins with same denom")
		} else {
			undertest := sdk.NewDecCoins(tt.input...)
			s.Require().Equal(tt.expected, undertest.String(), "NewDecCoins must return expected results")
			s.Require().Equal(tt.original, tt.input.String(), "input must be unmodified and match original")
		}
	}
}

func (s *decCoinTestSuite) TestIsValid() {
	tests := []struct {
		coin       sdk.DecCoin
		expectPass bool
		msg        string
	}{
		{
			sdk.NewDecCoin("mytoken", sdk.NewInt(10)),
			true,
			"valid coins should have passed",
		},
		{
			sdk.DecCoin{Denom: "BTC", Amount: sdk.NewDec(10)},
			true,
			"valid uppercase denom",
		},
		{
			sdk.DecCoin{Denom: "Bitcoin", Amount: sdk.NewDec(10)},
			true,
			"valid mixed case denom",
		},
		{
			sdk.DecCoin{Denom: "btc", Amount: sdk.NewDec(-10)},
			false,
			"negative amount",
		},
	}

	for _, tc := range tests {
		tc := tc
		if tc.expectPass {
			s.Require().True(tc.coin.IsValid(), tc.msg)
		} else {
			s.Require().False(tc.coin.IsValid(), tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestSubDecCoin() {
	tests := []struct {
		coin       sdk.DecCoin
		expectPass bool
		msg        string
	}{
		{
			sdk.NewDecCoin("mytoken", sdk.NewInt(20)),
			true,
			"valid coins should have passed",
		},
		{
			sdk.NewDecCoin("othertoken", sdk.NewInt(20)),
			false,
			"denom mismatch",
		},
		{
			sdk.NewDecCoin("mytoken", sdk.NewInt(9)),
			false,
			"negative amount",
		},
	}

	decCoin := sdk.NewDecCoin("mytoken", sdk.NewInt(10))

	for _, tc := range tests {
		tc := tc
		if tc.expectPass {
			equal := tc.coin.Sub(decCoin)
			s.Require().Equal(equal, decCoin, tc.msg)
		} else {
			s.Require().Panics(func() { tc.coin.Sub(decCoin) }, tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestSubDecCoins() {
	tests := []struct {
		coins      sdk.DecCoins
		expectPass bool
		msg        string
	}{
		{
			sdk.NewDecCoinsFromCoins(sdk.NewCoin("mytoken", sdk.NewInt(10)), sdk.NewCoin("btc", sdk.NewInt(20)), sdk.NewCoin("eth", sdk.NewInt(30))),
			true,
			"sorted coins should have passed",
		},
		{
			sdk.DecCoins{sdk.NewDecCoin("mytoken", sdk.NewInt(10)), sdk.NewDecCoin("btc", sdk.NewInt(20)), sdk.NewDecCoin("eth", sdk.NewInt(30))},
			false,
			"unorted coins should panic",
		},
		{
			sdk.DecCoins{sdk.DecCoin{Denom: "BTC", Amount: sdk.NewDec(10)}, sdk.NewDecCoin("eth", sdk.NewInt(15)), sdk.NewDecCoin("mytoken", sdk.NewInt(5))},
			false,
			"invalid denoms",
		},
	}

	decCoins := sdk.NewDecCoinsFromCoins(sdk.NewCoin("btc", sdk.NewInt(10)), sdk.NewCoin("eth", sdk.NewInt(15)), sdk.NewCoin("mytoken", sdk.NewInt(5)))

	for _, tc := range tests {
		tc := tc
		if tc.expectPass {
			equal := tc.coins.Sub(decCoins)
			s.Require().Equal(equal, decCoins, tc.msg)
		} else {
			s.Require().Panics(func() { tc.coins.Sub(decCoins) }, tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestSortDecCoins() {
	good := sdk.DecCoins{
		sdk.NewInt64DecCoin("gas", 1),
		sdk.NewInt64DecCoin("mineral", 1),
		sdk.NewInt64DecCoin("tree", 1),
	}
	empty := sdk.DecCoins{
		sdk.NewInt64DecCoin("gold", 0),
	}
	badSort1 := sdk.DecCoins{
		sdk.NewInt64DecCoin("tree", 1),
		sdk.NewInt64DecCoin("gas", 1),
		sdk.NewInt64DecCoin("mineral", 1),
	}
	badSort2 := sdk.DecCoins{ // both are after the first one, but the second and third are in the wrong order
		sdk.NewInt64DecCoin("gas", 1),
		sdk.NewInt64DecCoin("tree", 1),
		sdk.NewInt64DecCoin("mineral", 1),
	}
	badAmt := sdk.DecCoins{
		sdk.NewInt64DecCoin("gas", 1),
		sdk.NewInt64DecCoin("tree", 0),
		sdk.NewInt64DecCoin("mineral", 1),
	}
	dup := sdk.DecCoins{
		sdk.NewInt64DecCoin("gas", 1),
		sdk.NewInt64DecCoin("gas", 1),
		sdk.NewInt64DecCoin("mineral", 1),
	}
	cases := []struct {
		name          string
		coins         sdk.DecCoins
		before, after bool // valid before/after sort
	}{
		{"valid coins", good, true, true},
		{"empty coins", empty, false, false},
		{"unsorted coins (1)", badSort1, false, true},
		{"unsorted coins (2)", badSort2, false, true},
		{"zero amount coins", badAmt, false, false},
		{"duplicate coins", dup, false, false},
	}

	for _, tc := range cases {
		s.Require().Equal(tc.before, tc.coins.IsValid(), "coin validity is incorrect before sorting; %s", tc.name)
		tc.coins.Sort()
		s.Require().Equal(tc.after, tc.coins.IsValid(), "coin validity is incorrect after sorting;  %s", tc.name)
	}
}

func (s *decCoinTestSuite) TestDecCoinsValidate() {
	testCases := []struct {
		input        sdk.DecCoins
		expectedPass bool
	}{
		{sdk.DecCoins{}, true},
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}}, true},
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}, sdk.DecCoin{testDenom2, sdk.NewDec(100000)}}, true},
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(-5)}}, false},
		{sdk.DecCoins{sdk.DecCoin{"BTC", sdk.NewDec(5)}}, true},
		{sdk.DecCoins{sdk.DecCoin{"0BTC", sdk.NewDec(5)}}, false},
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}, sdk.DecCoin{"B", sdk.NewDec(100000)}}, false},
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}, sdk.DecCoin{testDenom2, sdk.NewDec(-100000)}}, false},
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(-5)}, sdk.DecCoin{testDenom2, sdk.NewDec(100000)}}, false},
		{sdk.DecCoins{sdk.DecCoin{"BTC", sdk.NewDec(5)}, sdk.DecCoin{testDenom2, sdk.NewDec(100000)}}, true},
		{sdk.DecCoins{sdk.DecCoin{"0BTC", sdk.NewDec(5)}, sdk.DecCoin{testDenom2, sdk.NewDec(100000)}}, false},
	}

	for i, tc := range testCases {
		err := tc.input.Validate()
		if tc.expectedPass {
			s.Require().NoError(err, "unexpected result for test case #%d, input: %v", i, tc.input)
		} else {
			s.Require().Error(err, "unexpected result for test case #%d, input: %v", i, tc.input)
		}
	}
}

func (s *decCoinTestSuite) TestParseDecCoins() {
	testCases := []struct {
		input          string
		expectedResult sdk.DecCoins
		expectedErr    bool
	}{
		{"", nil, false},
		{"4stake", sdk.DecCoins{sdk.NewDecCoinFromDec("stake", sdk.NewDecFromInt(sdk.NewInt(4)))}, false},
		{"5.5atom,4stake", sdk.DecCoins{
			sdk.NewDecCoinFromDec("atom", sdk.NewDecWithPrec(5500000000000000000, sdk.Precision)),
			sdk.NewDecCoinFromDec("stake", sdk.NewDec(4)),
		}, false},
		{"0.0stake", sdk.DecCoins{}, false}, // remove zero coins
		{"10.0btc,1.0atom,20.0btc", nil, true},
		{
			"0.004STAKE",
			sdk.DecCoins{sdk.NewDecCoinFromDec("STAKE", sdk.NewDecWithPrec(4000000000000000, sdk.Precision))},
			false,
		},
		{
			"0.004stake",
			sdk.DecCoins{sdk.NewDecCoinFromDec("stake", sdk.NewDecWithPrec(4000000000000000, sdk.Precision))},
			false,
		},
		{
			"5.04atom,0.004stake",
			sdk.DecCoins{
				sdk.NewDecCoinFromDec("atom", sdk.NewDecWithPrec(5040000000000000000, sdk.Precision)),
				sdk.NewDecCoinFromDec("stake", sdk.NewDecWithPrec(4000000000000000, sdk.Precision)),
			},
			false,
		},
		{"0.0stake,0.004stake,5.04atom", // remove zero coins
			sdk.DecCoins{
				sdk.NewDecCoinFromDec("atom", sdk.NewDecWithPrec(5040000000000000000, sdk.Precision)),
				sdk.NewDecCoinFromDec("stake", sdk.NewDecWithPrec(4000000000000000, sdk.Precision)),
			},
			false,
		},
	}

	for i, tc := range testCases {
		res, err := sdk.ParseDecCoins(tc.input)
		if tc.expectedErr {
			s.Require().Error(err, "expected error for test case #%d, input: %v", i, tc.input)
		} else {
			s.Require().NoError(err, "unexpected error for test case #%d, input: %v", i, tc.input)
			s.Require().Equal(tc.expectedResult, res, "unexpected result for test case #%d, input: %v", i, tc.input)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoinsString() {
	testCases := []struct {
		input    sdk.DecCoins
		expected string
	}{
		{sdk.DecCoins{}, ""},
		{
			sdk.DecCoins{
				sdk.NewDecCoinFromDec("atom", sdk.NewDecWithPrec(5040000000000000000, sdk.Precision)),
				sdk.NewDecCoinFromDec("stake", sdk.NewDecWithPrec(4000000000000000, sdk.Precision)),
			},
			"5.040000000000000000atom,0.004000000000000000stake",
		},
	}

	for i, tc := range testCases {
		out := tc.input.String()
		s.Require().Equal(tc.expected, out, "unexpected result for test case #%d, input: %v", i, tc.input)
	}
}

func (s *decCoinTestSuite) TestDecCoinsIntersect() {
	testCases := []struct {
		input1         string
		input2         string
		expectedResult string
	}{
		{"", "", ""},
		{"1.0stake", "", ""},
		{"1.0stake", "1.0stake", "1.0stake"},
		{"", "1.0stake", ""},
		{"1.0stake", "", ""},
		{"2.0stake,1.0trope", "1.9stake", "1.9stake"},
		{"2.0stake,1.0trope", "2.1stake", "2.0stake"},
		{"2.0stake,1.0trope", "0.9trope", "0.9trope"},
		{"2.0stake,1.0trope", "1.9stake,0.9trope", "1.9stake,0.9trope"},
		{"2.0stake,1.0trope", "1.9stake,0.9trope,20.0other", "1.9stake,0.9trope"},
		{"2.0stake,1.0trope", "1.0other", ""},
	}

	for i, tc := range testCases {
		in1, err := sdk.ParseDecCoins(tc.input1)
		s.Require().NoError(err, "unexpected parse error in %v", i)
		in2, err := sdk.ParseDecCoins(tc.input2)
		s.Require().NoError(err, "unexpected parse error in %v", i)
		exr, err := sdk.ParseDecCoins(tc.expectedResult)
		s.Require().NoError(err, "unexpected parse error in %v", i)
		s.Require().True(in1.Intersect(in2).IsEqual(exr), "in1.cap(in2) != exr in %v", i)
	}
}

func (s *decCoinTestSuite) TestDecCoinsTruncateDecimal() {
	decCoinA := sdk.NewDecCoinFromDec("bar", sdk.MustNewDecFromStr("5.41"))
	decCoinB := sdk.NewDecCoinFromDec("foo", sdk.MustNewDecFromStr("6.00"))

	testCases := []struct {
		input          sdk.DecCoins
		truncatedCoins sdk.Coins
		changeCoins    sdk.DecCoins
	}{
		{sdk.DecCoins{}, sdk.Coins(nil), sdk.DecCoins(nil)},
		{
			sdk.DecCoins{decCoinA, decCoinB},
			sdk.Coins{sdk.NewInt64Coin(decCoinA.Denom, 5), sdk.NewInt64Coin(decCoinB.Denom, 6)},
			sdk.DecCoins{sdk.NewDecCoinFromDec(decCoinA.Denom, sdk.MustNewDecFromStr("0.41"))},
		},
		{
			sdk.DecCoins{decCoinB},
			sdk.Coins{sdk.NewInt64Coin(decCoinB.Denom, 6)},
			sdk.DecCoins(nil),
		},
	}

	for i, tc := range testCases {
		truncatedCoins, changeCoins := tc.input.TruncateDecimal()
		s.Require().Equal(
			tc.truncatedCoins, truncatedCoins,
			"unexpected truncated coins; tc #%d, input: %s", i, tc.input,
		)
		s.Require().Equal(
			tc.changeCoins, changeCoins,
			"unexpected change coins; tc #%d, input: %s", i, tc.input,
		)
	}
}

func (s *decCoinTestSuite) TestDecCoinsQuoDecTruncate() {
	x := sdk.MustNewDecFromStr("1.00")
	y := sdk.MustNewDecFromStr("10000000000000000000.00")

	testCases := []struct {
		coins  sdk.DecCoins
		input  sdk.Dec
		result sdk.DecCoins
		panics bool
	}{
		{sdk.DecCoins{}, sdk.ZeroDec(), sdk.DecCoins(nil), true},
		{sdk.DecCoins{sdk.NewDecCoinFromDec("foo", x)}, y, sdk.DecCoins(nil), false},
		{sdk.DecCoins{sdk.NewInt64DecCoin("foo", 5)}, sdk.NewDec(2), sdk.DecCoins{sdk.NewDecCoinFromDec("foo", sdk.MustNewDecFromStr("2.5"))}, false},
	}

	for i, tc := range testCases {
		tc := tc
		if tc.panics {
			s.Require().Panics(func() { tc.coins.QuoDecTruncate(tc.input) })
		} else {
			res := tc.coins.QuoDecTruncate(tc.input)
			s.Require().Equal(tc.result, res, "unexpected result; tc #%d, coins: %s, input: %s", i, tc.coins, tc.input)
		}
	}
}

func (s *decCoinTestSuite) TestNewDecCoinsWithIsValid() {
	fake1 := append(sdk.NewDecCoins(sdk.NewDecCoin("mytoken", sdk.NewInt(10))), sdk.DecCoin{Denom: "10BTC", Amount: sdk.NewDec(10)})
	fake2 := append(sdk.NewDecCoins(sdk.NewDecCoin("mytoken", sdk.NewInt(10))), sdk.DecCoin{Denom: "BTC", Amount: sdk.NewDec(-10)})

	tests := []struct {
		coin       sdk.DecCoins
		expectPass bool
		msg        string
	}{
		{
			sdk.NewDecCoins(sdk.NewDecCoin("mytoken", sdk.NewInt(10))),
			true,
			"valid coins should have passed",
		},
		{
			fake1,
			false,
			"invalid denoms",
		},
		{
			fake2,
			false,
			"negative amount",
		},
	}

	for _, tc := range tests {
		tc := tc
		if tc.expectPass {
			s.Require().True(tc.coin.IsValid(), tc.msg)
		} else {
			s.Require().False(tc.coin.IsValid(), tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoins_AddDecCoinWithIsValid() {
	lengthTestDecCoins := sdk.NewDecCoins().Add(sdk.NewDecCoin("mytoken", sdk.NewInt(10))).Add(sdk.DecCoin{Denom: "BTC", Amount: sdk.NewDec(10)})
	s.Require().Equal(2, len(lengthTestDecCoins), "should be 2")

	tests := []struct {
		coin       sdk.DecCoins
		expectPass bool
		msg        string
	}{
		{
			sdk.NewDecCoins().Add(sdk.NewDecCoin("mytoken", sdk.NewInt(10))),
			true,
			"valid coins should have passed",
		},
		{
			sdk.NewDecCoins().Add(sdk.NewDecCoin("mytoken", sdk.NewInt(10))).Add(sdk.DecCoin{Denom: "0BTC", Amount: sdk.NewDec(10)}),
			false,
			"invalid denoms",
		},
		{
			sdk.NewDecCoins().Add(sdk.NewDecCoin("mytoken", sdk.NewInt(10))).Add(sdk.DecCoin{Denom: "BTC", Amount: sdk.NewDec(-10)}),
			false,
			"negative amount",
		},
	}

	for _, tc := range tests {
		tc := tc
		if tc.expectPass {
			s.Require().True(tc.coin.IsValid(), tc.msg)
		} else {
			s.Require().False(tc.coin.IsValid(), tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoins_Empty() {
	testCases := []struct {
		input          sdk.DecCoins
		expectedResult bool
		msg            string
	}{
		{sdk.DecCoins{}, true, "No coins as expected."},
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}}, false, "DecCoins is not empty"},
	}

	for _, tc := range testCases {
		if tc.expectedResult {
			s.Require().True(tc.input.Empty(), tc.msg)
		} else {
			s.Require().False(tc.input.Empty(), tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoins_GetDenomByIndex() {
	testCases := []struct {
		name           string
		input          sdk.DecCoins
		index          int
		expectedResult string
		expectedErr    bool
	}{
		{
			"No DecCoins in Slice",
			sdk.DecCoins{},
			0,
			"",
			true,
		},
		{"When index out of bounds", sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}}, 2, "", true},
		{"When negative index", sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}}, -1, "", true},
		{
			"Appropriate index case",
			sdk.DecCoins{
				sdk.DecCoin{testDenom1, sdk.NewDec(5)},
				sdk.DecCoin{testDenom2, sdk.NewDec(57)},
			}, 1, testDenom2, false,
		},
	}

	for i, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.expectedErr {
				s.Require().Panics(func() { tc.input.GetDenomByIndex(tc.index) }, "Test should have panicked")
			} else {
				res := tc.input.GetDenomByIndex(tc.index)
				s.Require().Equal(tc.expectedResult, res, "Unexpected result for test case #%d, expected output: %s, input: %v", i, tc.expectedResult, tc.input)
			}
		})
	}
}

func (s *decCoinTestSuite) TestDecCoins_IsAllPositive() {
	testCases := []struct {
		input          sdk.DecCoins
		expectedResult bool
		msg            string
	}{
		// No Coins
		{sdk.DecCoins{}, false, "No coins. Should be false."},

		// One Coin - Zero value
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(0)}}, false, "One coin with zero amount value. Should be false."},

		// One Coin - Postive value
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}}, true, "One coin with positive amount. Should be true."},

		// One Coin - Negative value
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(-15)}}, false, "One coin with negative amount. Should be false."},

		// Multiple Coins - All positive value
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(51)},
			sdk.DecCoin{testDenom1, sdk.NewDec(123)},
			sdk.DecCoin{testDenom1, sdk.NewDec(50)},
			sdk.DecCoin{testDenom1, sdk.NewDec(92233720)},
		}, true, "All positive amount. Should be true."},

		// Multiple Coins - Some negative value
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(51)},
			sdk.DecCoin{testDenom1, sdk.NewDec(-123)},
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
			sdk.DecCoin{testDenom1, sdk.NewDec(92233720)},
		}, false, "Not all positive amount. Should be false."},
	}

	for i, tc := range testCases {
		if tc.expectedResult {
			s.Require().True(tc.input.IsAllPositive(), "Test case #%d: %s", i, tc.msg)
		} else {
			s.Require().False(tc.input.IsAllPositive(), "Test case #%d: %s", i, tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoin_IsLT() {
	testCases := []struct {
		coin           sdk.DecCoin
		otherCoin      sdk.DecCoin
		expectedResult bool
		expectedPanic  bool
		msg            string
	}{

		// Same Denom - Less than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(3)}, sdk.DecCoin{testDenom1, sdk.NewDec(19)}, true, false, "DecCoin amount lesser than other coin. Should be true."},

		// Same Denom - Greater than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(343340)}, sdk.DecCoin{testDenom1, sdk.NewDec(14)}, false, false, "DecCoin amount greater than other coin. Should be false."},

		// Same Denom - Same as other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)}, sdk.DecCoin{testDenom1, sdk.NewDec(20)}, false, false, "DecCoin amount same as other coin. Should be false."},

		// Different Denom - Less than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(3)}, sdk.DecCoin{testDenom2, sdk.NewDec(19)}, true, true, "DecCoin denom different than other coin. Should panic."},

		// Different Denom - Greater than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(343340)}, sdk.DecCoin{testDenom2, sdk.NewDec(14)}, true, true, "DecCoin denom different than other coin. Should panic."},

		// Different Denom - Same as other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)}, sdk.DecCoin{testDenom2, sdk.NewDec(20)}, true, true, "DecCoin denom different than other coin. Should panic."},
	}

	for i, tc := range testCases {
		if tc.expectedPanic {
			s.Require().Panics(func() { tc.coin.IsLT(tc.otherCoin) }, "Test case #%d: %s", i, tc.msg)
		} else {
			res := tc.coin.IsLT(tc.otherCoin)
			if tc.expectedResult {
				s.Require().True(res, "Test case #%d: %s", i, tc.msg)
			} else {
				s.Require().False(res, "Test case #%d: %s", i, tc.msg)
			}
		}
	}
}

func (s *decCoinTestSuite) TestDecCoin_IsGTE() {
	testCases := []struct {
		coin           sdk.DecCoin
		otherCoin      sdk.DecCoin
		expectedResult bool
		expectedPanic  bool
		msg            string
	}{

		// Same Denom - Less than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(3)}, sdk.DecCoin{testDenom1, sdk.NewDec(19)}, false, false, "DecCoin amount lesser than other coin. Should be false."},

		// Same Denom - Greater than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(343340)}, sdk.DecCoin{testDenom1, sdk.NewDec(14)}, true, false, "DecCoin amount greater than other coin. Should be true."},

		// Same Denom - Same as other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)}, sdk.DecCoin{testDenom1, sdk.NewDec(20)}, true, false, "DecCoin amount equal to other coin. Should be true."},

		// Different Denom - Less than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(3)}, sdk.DecCoin{testDenom2, sdk.NewDec(19)}, true, true, "DecCoin denom different than other coin. Should panic."},

		// Different Denom - Greater than other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(343340)}, sdk.DecCoin{testDenom2, sdk.NewDec(14)}, true, true, "DecCoin denom different than other coin. Should panic."},

		// Different Denom - Same as other coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)}, sdk.DecCoin{testDenom2, sdk.NewDec(20)}, true, true, "DecCoin denom different than other coin. Should panic."},
	}

	for i, tc := range testCases {
		if tc.expectedPanic {
			s.Require().Panics(func() { tc.coin.IsGTE(tc.otherCoin) }, "Test case #%d: %s", i, tc.msg)
		} else {
			res := tc.coin.IsGTE(tc.otherCoin)
			if tc.expectedResult {
				s.Require().True(res, "Test case #%d: %s", i, tc.msg)
			} else {
				s.Require().False(res, "Test case #%d: %s", i, tc.msg)
			}
		}
	}
}

func (s *decCoinTestSuite) TestDecCoins_IsZero() {
	testCases := []struct {
		coins          sdk.DecCoins
		expectedResult bool
		msg            string
	}{
		// No Coins
		{sdk.DecCoins{}, true, "No coins. Should be true."},

		// One Coin - Zero value
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(0)}}, true, "One coin with zero amount value. Should be true."},

		// One Coin - Postive value
		{sdk.DecCoins{sdk.DecCoin{testDenom1, sdk.NewDec(5)}}, false, "One coin with positive amount. Should be false."},

		// Multiple Coins - All zero value
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
		}, true, "All zero amount coins. Should be true."},

		// Multiple Coins - Some positive value
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
			sdk.DecCoin{testDenom1, sdk.NewDec(0)},
			sdk.DecCoin{testDenom1, sdk.NewDec(92233720)},
		}, false, "Not all zero amount coins. Should be false."},
	}

	for i, tc := range testCases {
		if tc.expectedResult {
			s.Require().True(tc.coins.IsZero(), "Test case #%d: %s", i, tc.msg)
		} else {
			s.Require().False(tc.coins.IsZero(), "Test case #%d: %s", i, tc.msg)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoins_MulDec() {
	testCases := []struct {
		coins          sdk.DecCoins
		multiplier     sdk.Dec
		expectedResult sdk.DecCoins
		msg            string
	}{
		// No Coins
		{sdk.DecCoins{}, sdk.NewDec(1), sdk.DecCoins(nil), "No coins. Should return empty slice."},

		// Multiple coins - zero multiplier
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(10)},
			sdk.DecCoin{testDenom1, sdk.NewDec(30)},
		}, sdk.NewDec(0), sdk.DecCoins(nil), "Multipler is zero. Should return empty slice."},

		// Multiple coins - positive multiplier
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(1)},
			sdk.DecCoin{testDenom1, sdk.NewDec(2)},
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom1, sdk.NewDec(4)},
		}, sdk.NewDec(2), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(20)},
		}, "Multipler is positive. Should return multiplied deccoins."},

		// Multiple coins - negative multiplier
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(1)},
			sdk.DecCoin{testDenom1, sdk.NewDec(2)},
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom1, sdk.NewDec(4)},
		}, sdk.NewDec(-2), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(-20)},
		}, "Multipler is negative. Should return multiplied deccoins."},

		// Multiple coins - Different denom
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(1)},
			sdk.DecCoin{testDenom2, sdk.NewDec(2)},
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom2, sdk.NewDec(4)},
		}, sdk.NewDec(2), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(8)},
			sdk.DecCoin{testDenom2, sdk.NewDec(12)},
		}, "Multiple coins with different denoms. Should return multiplied deccoins with appropriate denoms."},
	}

	for i, tc := range testCases {
		res := tc.coins.MulDec(tc.multiplier)
		s.Require().Equal(tc.expectedResult, res, "Test case #%d: %s %s", i, tc.msg, res)
	}
}

func (s *decCoinTestSuite) TestDecCoins_MulDecTruncate() {
	testCases := []struct {
		coins          sdk.DecCoins
		multiplier     sdk.Dec
		expectedResult sdk.DecCoins
		expectedPanic  bool
		msg            string
	}{
		// No Coins
		{sdk.DecCoins{}, sdk.NewDec(1), sdk.DecCoins(nil), false, "No coins. Should return nil."},

		// Multiple coins - zero multiplier
		// TODO - Fix test - Function comment documentation for MulDecTruncate says if multiplier d is zero, it should panic.
		// However, that is not the observed behaviour. Currently nil is returned.
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(10, 3)},
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(30, 2)},
		}, sdk.NewDec(0), sdk.DecCoins(nil), false, "Multipler is zero. Should return nil."},

		// Multiple coins - positive multiplier
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(15, 1)},
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(15, 1)},
		}, sdk.NewDec(1), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(3, 0)},
		}, false, "Multipler is positive. Should return truncated multiplied deccoins."},

		// Multiple coins - positive multiplier
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(15, 1)},
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(15, 1)},
		}, sdk.NewDec(-2), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(-6, 0)},
		}, false, "Multipler is positive. Should return truncated multiplied deccoins."},

		// Multiple coins - Different denom
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(15, 1)},
			sdk.DecCoin{testDenom2, sdk.NewDecWithPrec(3333, 4)},
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(15, 1)},
			sdk.DecCoin{testDenom2, sdk.NewDecWithPrec(333, 4)},
		}, sdk.NewDec(10), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(30, 0)},
			sdk.DecCoin{testDenom2, sdk.NewDecWithPrec(3666, 3)},
		}, false, "Multipler is positive. Should return truncated multiplied deccoins with appropriate denoms."},
	}

	for i, tc := range testCases {
		if tc.expectedPanic {
			s.Require().Panics(func() { tc.coins.MulDecTruncate(tc.multiplier) }, "Test case #%d: %s", i, tc.msg)
		} else {
			res := tc.coins.MulDecTruncate(tc.multiplier)
			s.Require().Equal(tc.expectedResult, res, "Test case #%d: %s %s", i, tc.msg, res)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoins_QuoDec() {

	testCases := []struct {
		coins          sdk.DecCoins
		input          sdk.Dec
		expectedResult sdk.DecCoins
		panics         bool
		msg            string
	}{
		// No Coins
		{sdk.DecCoins{}, sdk.NewDec(1), sdk.DecCoins(nil), false, "No coins. Should return empty slice."},

		// Multiple coins - zero input
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(10)},
			sdk.DecCoin{testDenom1, sdk.NewDec(30)},
		}, sdk.NewDec(0), sdk.DecCoins(nil), true, "Input is zero. Should panic."},

		// Multiple coins - positive input
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom1, sdk.NewDec(4)},
		}, sdk.NewDec(2), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(35, 1)},
		}, false, "Input is positive. Should return divided deccoins."},

		// Multiple coins - negative input
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom1, sdk.NewDec(4)},
		}, sdk.NewDec(-2), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDecWithPrec(-35, 1)},
		}, false, "Input is negative. Should return divided deccoins."},

		// Multiple coins - Different input
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(1)},
			sdk.DecCoin{testDenom2, sdk.NewDec(2)},
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom2, sdk.NewDec(4)},
		}, sdk.NewDec(2), sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(2)},
			sdk.DecCoin{testDenom2, sdk.NewDec(3)},
		}, false, "Input coins with different denoms. Should return divided deccoins with appropriate denoms."},
	}

	for i, tc := range testCases {
		tc := tc
		if tc.panics {
			s.Require().Panics(func() { tc.coins.QuoDec(tc.input) })
		} else {
			res := tc.coins.QuoDec(tc.input)
			s.Require().Equal(tc.expectedResult, res, "Test case #%d: %s %s", i, tc.msg, res)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoin_IsEqual() {
	testCases := []struct {
		coin           sdk.DecCoin
		otherCoin      sdk.DecCoin
		expectedResult bool
		expectedPanic  bool
		msg            string
	}{

		// Different Denom Same Amount
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)},
			sdk.DecCoin{testDenom2, sdk.NewDec(20)},
			false, true, "Different demon for coins. Should panic."},

		// Different Denom Different Amount
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)},
			sdk.DecCoin{testDenom2, sdk.NewDec(10)},
			false, true, "Different demon for coins.. Should panic."},

		// Same Denom Different Amount
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)},
			sdk.DecCoin{testDenom1, sdk.NewDec(10)},
			false, false, "Same denom but different amount. Should be false."},

		// Same Denom Same Amount
		{sdk.DecCoin{testDenom1, sdk.NewDec(20)},
			sdk.DecCoin{testDenom1, sdk.NewDec(20)},
			true, false, "Same denom and same amount. Should be true."},
	}

	for i, tc := range testCases {
		if tc.expectedPanic {
			s.Require().Panics(func() { tc.coin.IsEqual(tc.otherCoin) }, "Test case #%d: %s", i, tc.msg)
		} else {
			res := tc.coin.IsEqual(tc.otherCoin)
			if tc.expectedResult {
				s.Require().True(res, "Test case #%d: %s", i, tc.msg)
			} else {
				s.Require().False(res, "Test case #%d: %s", i, tc.msg)
			}
		}
	}
}

func (s *decCoinTestSuite) TestDecCoins_IsEqual() {
	testCases := []struct {
		coinsA         sdk.DecCoins
		coinsB         sdk.DecCoins
		expectedResult bool
		expectedPanic  bool
		msg            string
	}{
		// Different length sets
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom1, sdk.NewDec(4)},
		}, sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(35)},
		}, false, false, "Different length coin sets. Should be false."},

		// Same length - different denoms
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom1, sdk.NewDec(4)},
		}, sdk.DecCoins{
			sdk.DecCoin{testDenom2, sdk.NewDec(3)},
			sdk.DecCoin{testDenom2, sdk.NewDec(4)},
		}, false, true, "Same length coin sets with different denoms. Should panic."},

		// Same length - different amounts
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(3)},
			sdk.DecCoin{testDenom1, sdk.NewDec(4)},
		}, sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(41)},
			sdk.DecCoin{testDenom1, sdk.NewDec(343)},
		}, false, false, "Same length coin sets with different amounts. Should be false."},

		// Same length - same amounts
		{sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(33)},
			sdk.DecCoin{testDenom1, sdk.NewDec(344)},
		}, sdk.DecCoins{
			sdk.DecCoin{testDenom1, sdk.NewDec(33)},
			sdk.DecCoin{testDenom1, sdk.NewDec(344)},
		}, true, false, "Same length coin sets with same amounts and denoms. Should be true,"},
	}

	for i, tc := range testCases {
		if tc.expectedPanic {
			s.Require().Panics(func() { tc.coinsA.IsEqual(tc.coinsB) }, "Test case #%d: %s", i, tc.msg)
		} else {
			res := tc.coinsA.IsEqual(tc.coinsB)
			if tc.expectedResult {
				s.Require().True(res, "Test case #%d: %s", i, tc.msg)
			} else {
				s.Require().False(res, "Test case #%d: %s", i, tc.msg)
			}
		}
	}
}

func (s *decCoinTestSuite) TestDecCoin_Validate() {
	var empty sdk.DecCoin
	testCases := []struct {
		input        sdk.DecCoin
		expectedPass bool
	}{
		// uninitalized deccoin
		{empty, false},

		// invalid denom string
		{sdk.DecCoin{"(){9**&})", sdk.NewDec(33)}, false},

		// negative coin amount
		{sdk.DecCoin{testDenom1, sdk.NewDec(-33)}, false},

		// valid coin
		{sdk.DecCoin{testDenom1, sdk.NewDec(33)}, true},
	}

	for i, tc := range testCases {
		err := tc.input.Validate()
		if tc.expectedPass {
			s.Require().NoError(err, "unexpected result for test case #%d, input: %v", i, tc.input)
		} else {
			s.Require().Error(err, "unexpected result for test case #%d, input: %v", i, tc.input)
		}
	}
}

func (s *decCoinTestSuite) TestDecCoin_ParseDecCoin() {
	var empty sdk.DecCoin
	testCases := []struct {
		input          string
		expectedResult sdk.DecCoin
		expectedErr    bool
	}{
		// empty input
		{"", empty, true},

		// bad input
		{"✨🌟⭐", empty, true},

		// invalid decimal coin
		{"9.3.0stake", empty, true},

		// precision over limit
		{"9.11111111111111111111stake", empty, true},

		// invalid denom
		// TODO - Clarify - According to error message for ValidateDenom call, supposed to
		// throw error when upper case characters are used. Currently uppercase denoms are allowed.
		{"9.3STAKE", sdk.DecCoin{"STAKE", sdk.NewDecWithPrec(93, 1)}, false},

		// valid input - amount and denom seperated by space
		{"9.3 stake", sdk.DecCoin{"stake", sdk.NewDecWithPrec(93, 1)}, false},

		// valid input - amount and denom concatenated
		{"9.3stake", sdk.DecCoin{"stake", sdk.NewDecWithPrec(93, 1)}, false},
	}

	for i, tc := range testCases {
		res, err := sdk.ParseDecCoin(tc.input)
		if tc.expectedErr {
			s.Require().Error(err, "expected error for test case #%d, input: %v", i, tc.input)
		} else {
			s.Require().NoError(err, "unexpected error for test case #%d, input: %v", i, tc.input)
			s.Require().Equal(tc.expectedResult, res, "unexpected result for test case #%d, input: %v", i, tc.input)
		}
	}
}
