/*
Copyright (c) 2023, Balazs Nyiro
All rights reserved.

This source code (all file in this repo) is licensed
under the Apache-2 style license found in the
LICENSE file in the root directory of this source tree.
*/
package erlango

import (
	"fmt"
	"testing"
)

//  go test -v -run Test_numbers_integers
func Test_numbers_integers(t *testing.T) {
	funName := "Test_numbers_integers"
	fmt.Println(funName)


	erlSrc :=`	VarAtom = atomValue1, 
				VarAtomQuoted = 'atomQuoted', 
				Str = "str", 
				Float = 1.1, 


				TestGoal = "fullNumberTestTable"

				IntSimple = 12, 
				IntSimpleUnderscored = 12_34, 
				IntegerWithUnderscore = 0_13.

				DollarChar = $A,
				DollarWithTwoChars = $\n,

			
				BaseValueSimple = 2#101,

				BaseValueWithUnderscore = 2#1_01,


				BaseValueHexaLowerCap = 16#1f,
				BaseValueHexaLowerCap_ff = 16#ff,
				Comment1 = "ff can be detected as an atom, in prev steps"

				BaseValueHexaUpperCap = 16#1F,
				BaseValueHexaUpperCap_FF = 16#FF,
				Comment2 = "FF can be detected, as a variable name"

		



				BaseValueHexaWithUnderscoreInBase1 = 1_6#1f,
				BaseValueHexaWithUnderscoreInBase2 = 1_6#2F,
				BaseValueHexaWithUnderscoreEverywhere = 1_6#a_b

				FloatSimple = 12.34,
				FloatSimpleUnderscoreInInts = 12_34.56,
				FloatSimpleUnderscoreEverywhere = 12_34.56_78,

				ScientificSimple = 2.3e3,
				ScientificPlus = 2.3e+3,
				ScientificMinus = 2.3e-3,
				ScientificUnderscorePlus1 = 2_3.4e+3,

				ScientificUnderscoreEverywherePlus = 2_3.4e+3_0,
				IllegalFloatError_notValidErlangNumberIfSpaceInserted = "2_3.4e +3_0"

				ScientificCapical = 2.0E3,
				

				HexaCrazyLong = 16#11111111111111111111111111111,
				HexaCrazyLongDecimalValue = 5538449982437149470432529417834769,


				CommentEnd1 = "scientific notation can be mixed with non-decimal numbers: (val 26)",
				ScientificHexa = 16#1e-4.


`

	wantedExpressionDetectionTypes := "detectAllExpressions"
	erlExpressions := Expression_detection_for_tests(erlSrc, wantedExpressionDetectionTypes)
	erlExpressions.printAll()

}
