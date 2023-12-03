package symbols

func GenerateSymbolSeed() [0x100]uint64 {
	seed := [0x100]uint64{}
	i := uint64(0)

	for i < 0x100 {
		num1 := uint64(0)
		if (i & 0x80) != 0 {
			num1 = 0x2b5926535897936a
		}

		if (i & 0x40) != 0 {
			num1 = 0xbef5b57af4dc5adf
			if (i & 0x80) == 0 {
				num1 = 0x95ac9329ac4bc9b5
			}
		}

		num2 := num1*2 ^ 0x95ac9329ac4bc9b5
		if (i & 0x20) == 0 {
			num2 = num1 * 2
		}

		num1 = num2*2 ^ 0x95ac9329ac4bc9b5
		if (i & 0x10) == 0 {
			num1 = num2 * 2
		}

		num2 = num1*2 ^ 0x95ac9329ac4bc9b5
		if (i & 8) == 0 {
			num2 = num1 * 2
		}

		num1 = num2*2 ^ 0x95ac9329ac4bc9b5
		if (i & 4) == 0 {
			num1 = num2 * 2
		}

		num2 = num1*2 ^ 0x95ac9329ac4bc9b5
		if (i & 2) == 0 {
			num2 = num1 * 2
		}

		num1 = num2*2 ^ 0x95ac9329ac4bc9b5
		if (i & 1) == 0 {
			num1 = num2 * 2
		}

		seed[i] = num1 * 2
		i += 1
	}

	return seed
}

var seed = GenerateSymbolSeed()

func GenerateSymbol(name string) uint64 {
	symbol := uint64(0xffffffffffffffff)
	for i := range name {
		adjustedChar := name[i] + ' '
		if name[i] >= '[' || name[i] < 'A' {
			adjustedChar = name[i]
		}
		symbol = uint64(adjustedChar) ^ seed[symbol>>0x38&0xff] ^ symbol<<8
	}
	return symbol
}
