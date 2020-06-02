// ---------------------------- 35 original CHIP8 opcodes ---------------------------- //
// 0NNN: Execute RCA 1802 machine language routine at address NNN
// 00E0: Clear the screen
// 00EE: Return from subroutine
// 1NNN: Jump to address NNN
// 2NNN: Call subroutine at address NNN
// 3XNN: Skip the following instruction if the value of register VX equals NN
// 4XNN: Skip the following instruction if the value of register VX is not equal to NN
// 5XY0: Skip the following instruction if the value of register VX is equal to the value of register VY
// 6XNN: Set VX to NN
// 7XNN: Add NN to VX
// 8XY0: Set VX to the value in VY
// 8XY1: Set VX to VX OR VY
// 8XY2: Set VX to VX AND VY
// 8XY3: Set VX to VX XOR VY
// 8XY4: Add the value of register VY to register VX. Set VF to 01 if a carry occurs. Set VF to 00 if a carry does not occur
// 8XY5: Subtract the value of register VY from register VX. Set VF to 00 if a borrow occurs. Set VF to 01 if a borrow does not occur
// 8XY6: Store the value of register VY shifted right one bit in register VX. Set register VF to the least significant bit prior to the shift
// 8XY7: Set register VX to the value of VY minus VX. Set VF to 00 if a borrow occurs. Set VF to 01 if a borrow does not occur
// 8XYE: Store the value of register VY shifted left one bit in register VX. Set register VF to the most significant bit prior to the shift
// 9XY0: Skip the following instruction if the value of register VX is not equal to the value of register VY
// ANNN: Store memory address NNN in register I
// BNNN: Jump to address NNN + V0
// CXNN: Set VX to a random number with a mask of NN
// DXYN: Draw a sprite at position VX, VY with N bytes of sprite data starting at the address stored in I. Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
// EX9E: Skip the following instruction if the key corresponding to the hex value currently stored in register VX is pressed
// EXA1: Skip the following instruction if the key corresponding to the hex value currently stored in register VX is not pressed
// FX07: Store the current value of the delay timer in register VX
// FX0A: Wait for a keypress and store the result in register VX
// FX15: Set the delay timer to the value of register VX
// FX18: Set the sound timer to the value of register VX
// FX1E: Add the value stored in register VX to register I
// FX29: Set I to the memory address of the sprite data corresponding to the hexadecimal digit stored in register VX
// FX33: Store the binary-coded decimal equivalent of the value stored in register VX at addresses I, I+1, and I+2
// FX55: Store the values of registers V0 to VX inclusive in memory starting at address I. I is set to I + X + 1 after operation
// FX65: Fill registers V0 to VX inclusive with the values stored in memory starting at address I. I is set to I + X + 1 after operation
// --------------------- 02D8 undocumented opcode (maybe 0NNN 1802 instrunction?) -------------------- //
// 02D8: LDA 02, I // Load from memory at address I into V[00] to V[02]

package CPU

import (
	"os"
	"fmt"
	"math/rand"
)

// ---------------------------- CHIP-8 0xxx instruction set ---------------------------- //

// 0NNN
// Execute RCA 1802 machine language routine at address NNN
func opc_chip8_0NNN() {
	// Not needed by any game, just for documentation
}

// 00E0 - CLS
// Clear the display.
func opc_chip8_00E0() {
	Graphics = [128 * 64]byte{}
	PC += 2
	if Debug {
		fmt.Println("\t\tOpcode 00E0 executed. - Clear the display\n")
	}
}

// 00EE - RET
// Return from a subroutine
// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
// Need to incremente PC (PC+=2) after receive the value from Stack
func opc_chip8_00EE() {
	PC = Stack[SP] + 2
	SP --
	if Debug {
		fmt.Printf("\t\tOpcode 00EE executed. - Return from a subroutine (PC=%d)\n", PC)
	}
}

// ---------------------------- CHIP-8 1xxx instruction set ---------------------------- //

// 1nnn - JP addr
// Jump to location nnn.
// The interpreter sets the program counter to nnn.
func opc_chip8_1NNN() {
	PC = Opcode & 0x0FFF
	if Debug {
		fmt.Printf("\t\tOpcode 1nnn executed: Jump to location 0x%d\n", Opcode & 0x0FFF)
	}
}

// ---------------------------- CHIP-8 2xxx instruction set ---------------------------- //

// 2nnn - CALL addr
// Call subroutine at nnn.
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
func opc_chip8_2NNN() {
	SP++
	Stack[SP] = PC
	PC = Opcode & 0x0FFF
	if Debug {
		fmt.Printf("\t\tOpcode 2nnn executed: Call Subroutine at 0x%d\n", PC)
	}
}

// ---------------------------- CHIP-8 3xxx instruction set ---------------------------- //

// 3xnn - SE Vx, byte
// Skip next instruction if Vx = NN.
// The interpreter compares register Vx to nn, and if they are equal, increments the program counter by 2.
func opc_chip8_3XNN() {
	x := (Opcode & 0x0F00) >> 8
	nn := byte(Opcode & 0x00FF)
	if V[x] == nn {
		PC += 4
		if Debug {
			fmt.Printf("\t\tOpcode 3xnn executed: V[x(%d)]:(%d) = nn(%d), skip one instruction.\n", x, V[x], nn)
		}
	} else {
		PC += 2
		if Debug {
			fmt.Printf("\t\tOpcode 3xnn executed: V[x(%d)]:(%d) != nn(%d), do NOT skip one instruction.\n", x, V[x], nn)
		}
	}
}

// ---------------------------- CHIP-8 4xxx instruction set ---------------------------- //

// 4xnn - SNE Vx, byte
// Skip next instruction if Vx != nn.
// The interpreter compares register Vx to nn, and if they are not equal, increments the program counter by 2.
func opc_chip8_4XNN() {
	x := (Opcode & 0x0F00) >> 8
	nn := byte(Opcode & 0x00FF)
	if V[x] != nn {
		PC += 4
		if Debug {
			fmt.Printf("\t\tOpcode 4xnn executed: V[x(%d)]:%d != nn(%d), skip one instruction\n", x, V[x], nn)
		}
	} else {
		if Debug {
			fmt.Printf("\t\tOpcode 4xnn executed: V[x(%d)]:%d = nn(%d), NOT skip one instruction\n", x, V[x], nn)
		}
		PC += 2
	}
}

// ---------------------------- CHIP-8 5xxx instruction set ---------------------------- //

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
func opc_chip8_5XY0() {
	x := (Opcode & 0x0F00) >> 8
	y := (Opcode & 0x00F0) >> 4

	if (V[x] == V[y]){
		PC += 4
		if Debug {
			fmt.Printf("\t\tOpcode 5xy0 executed: V[x(%d)]:%d EQUAL V[y(%d)]:%d, SKIP one instruction\n", x, V[x], y, V[y])
		}
	} else {
		PC += 2
		if Debug {
			fmt.Printf("\t\tOpcode 5xy0 executed: V[x(%d)]:%d NOT EQUAL V[y(%d)]:%d, DO NOT SKIP one instruction\n", x, V[x], y, V[y])
		}
	}
}

// ---------------------------- CHIP-8 6xxx instruction set ---------------------------- //

// 6xnn - LD Vx, byte
// Set Vx = nn.
// The interpreter puts the value nn into register Vx.
func opc_chip8_6XNN() {
	x := (Opcode & 0x0F00) >> 8
	nn := byte(Opcode)

	V[x] = nn
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 6xnn executed: Set V[x(%d)] = %d\n", x, nn)
	}
}

// ---------------------------- CHIP-8 7xxx instruction set ---------------------------- //

// 7xnn - ADD Vx, byte
// Set Vx = Vx + nn.
// Adds the value nn to the value of register Vx, then stores the result in Vx.
func opc_chip8_7XNN() {
	x := (Opcode & 0x0F00) >> 8
	nn := byte(Opcode)

	V[x] += nn

	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 7xnn executed: Add the value nn(%d) to V[x(%d)]\n", nn, x)
	}
}


// ---------------------------- CHIP-8 8xxx instruction set ---------------------------- //

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
// Stores the value of register Vy in register Vx.
func opc_chip8_8XY0(x, y uint16) {
	V[x] = V[y]
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 8xy0 executed: Set V[x(%d)] = V[y(%d)]:%d\n", x, y, V[y])
	}
}

// 8xy1 - Set Vx = Vx OR Vy.
// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A bitwise OR compares the corrseponding bits from two values,
// and if either bit is 1, then the same bit in the result is also 1. Otherwise, it is 0.
func opc_chip8_8XY1(x, y uint16) {
	V[x] |= V[y]
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 8xy1 executed: Set V[x(%d)]:%d OR V[y(%d)]:%d\n", x, V[x], y, V[y])
	}
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, and if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
func opc_chip8_8XY2(x, y uint16) {
	V[x] &= V[y]
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 8xy2 executed: Set V[x(%d)] = V[x(%d)] AND V[y(%d)]\n", x, x, y)
	}
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values,
// and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func opc_chip8_8XY3(x, y uint16) {
	if Debug {
		fmt.Printf("\t\tOpcode 8xy3 executed:  V[x(%d)]:%d XOR V[y(%d)]:%d \n", x, V[x], y, V[y])
	}
	V[x] ^= V[y]
	PC += 2
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0.
// Only the lowest 8 bits of the result are kept, and stored in Vx.
func opc_chip8_8XY4(x, y uint16) {
	if ( V[x] + V[y] < V[x]) {
		V[0xF] = 1
	} else {
		V[0xF] = 0
	}
	if Debug {
		fmt.Printf("\t\tOpcode 8xy4 executed: Set V[x(%d)] = V[x(%d)] + V[y(%d)]\n", x, x, y)
	}
	// Old implementation, sum values, READ THE DOCS IN CASE OF PROBLEMS
	V[x] += V[y]

	PC += 2
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
func opc_chip8_8XY5(x, y uint16) {
	if V[x] >= V[y] {
		V[0xF] = 1
	} else {
		V[0xF] = 0
	}

	V[x] -= V[y]
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 8xy5 executed: Set V[x(%d)] = V[x(%d)]:%d - V[y(%d)]:%d\n", x, x, V[x], y, V[y])
	}
}

// 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2 (SHR).
// Original Chip8 INCREMENT I in this instruction
func opc_chip8_8XY6(x, y uint16) {
	V[0xF] = V[x] & 0x01

	if Legacy_8xy6_8xyE {
		V[x] = V[y] >> 1
	} else {
		V[x] = V[x] >> 1
	}

	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 8xy6 executed: Set V[x(%d)]:%d SHIFT RIGHT 1\n", x, V[x])
	}
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func opc_chip8_8XY7(x, y uint16) {
	if V[x] > V[y] {
		V[0xF] = 0
	} else {
		V[0xF] = 1
	}
	if Debug {
		fmt.Printf("\t\tOpcode 8xy7 executed: Set V[x(%d)]:%d = V[y(%d)]:%d - V[x(%d)]:%d\t\t = %d \n", x, V[x], y, V[y], x, V[x], V[y] - V[x])
	}
	V[x] = V[y] - V[x]

	PC += 2
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
func opc_chip8_8XYE(x, y uint16) {
	V[0xF] = V[x] >> 7 // Set V[F] to the Most Important Bit

	if Legacy_8xy6_8xyE {
		V[x] = V[y] << 1
	} else {
		V[x] = V[x] << 1
	}

	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 8xyE executed: Set V[x(%d)]:%d SHIFT LEFT 1\n", x, V[x])
	}
}

// ---------------------------- CHIP-8 9xxx instruction set ---------------------------- //

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
func opc_chip8_9XY0() {
	x := (Opcode & 0x0F00) >> 8
	y := (Opcode & 0x00F0) >> 4

	if ( V[x] != V[y] ) {
		PC += 4
		if Debug {
			fmt.Printf("\t\tOpcode 9xy0 executed: V[x(%d)]:%d != V[y(%d)]:%d, SKIP one instruction\n", x, V[x], y, V[y])
		}
	} else {
		PC += 2
		if Debug {
			fmt.Printf("\t\tOpcode 9xy0 executed: V[x(%d)]:%d = V[y(%d)]:%d, DO NOT SKIP one instruction\n", x, V[x], y, V[y])
		}
	}
}

// ---------------------------- CHIP-8 Axxx instruction set ---------------------------- //

// Annn - LD I, addr
// Set I = nnn.
// The value of register I is set to nnn.
func opc_chip8_ANNN() {
	I = Opcode & 0x0FFF
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode Annn executed: Set I = %d\n", I)
	}
}

// ---------------------------- CHIP-8 Bxxx instruction set ---------------------------- //

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
// The program counter is set to nnn plus the value of V0.
func opc_chip8_BNNN() {
	nnn := Opcode & 0x0FFF
	PC = nnn + uint16(V[0])
	if Debug {
		print ("\t\tOpcode Bnnn executed: Jump to location nnn(%d) + V[0(%d)]\n", nnn, V[0])
	}
}

// ---------------------------- CHIP-8 Cxxx instruction set ---------------------------- //

// Cxnn - RND Vx, byte
// Set Vx = random byte AND nn.
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value nn. The results are stored in Vx. See instruction 8xy2 for more information on AND.
func opc_chip8_CXNN() {
	x := uint16(Opcode&0x0F00) >> 8
	nn := Opcode & 0x00FF
	V[x] = byte(rand.Float32()*255) & byte(nn)
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode Cxnn executed: V[x(%d)] = %d (random byte AND nn(%d)) = %d\n", x, V[x], nn, V[x])
	}
}

// ---------------------------- CHIP-8 Dxxx instruction set ---------------------------- //

	//LATER

// ---------------------------- CHIP-8 Exxx instruction set ---------------------------- //

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
func opc_chip8_EX9E(x uint16) {
	if Key[V[x]] == 1 {
		PC += 4
		if Debug {
			fmt.Printf("\t\tOpcode Ex9E executed: Key[%d] pressed, skip one instruction\n",V[x])
		}
	} else {
		PC += 2
		if Debug {
			fmt.Printf("\t\tOpcode Ex9E executed: Key[%d] NOT pressed, continue\n",V[x])
		}
	}
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
func opc_chip8_EXA1(x uint16) {
	if Key[V[x]] == 0 {
		PC += 4
		if Debug {
			fmt.Printf("\t\tOpcode ExA1 executed: Key[%d] NOT pressed, skip one instruction\n",V[x])
		}
	} else {
		Key[V[x]] = 0
		PC += 2
		if Debug {
			fmt.Printf("\t\tOpcode ExA1 executed: Key[%d] pressed, continue\n",V[x])
		}
	}
}

// ---------------------------- CHIP-8 Fxxx instruction set ---------------------------- //

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
// The value of DT is placed into Vx.
func opc_chip8_FX07(x uint16) {
	V[x] = DelayTimer
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode Fx07 executed: Set V[x(%d)] with value of DelayTimer(%d)\n", x, DelayTimer)
	}
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func opc_chip8_FX0A(x uint16) {
	pressed := 0
	for i := 0 ; i < len(Key) ; i++ {
		if (Key[i] == 1){
			V[x] = byte(i)
			pressed = 1
			PC +=2
			if Debug {
				fmt.Printf("\t\tOpcode Fx0A executed: Wait for a key (Key[%d]) press -  (PRESSED)\n", i)
			}
			// Stop after find the first key pressed
			break
		}
	}
	if pressed == 0 {
		if Debug {
			fmt.Printf("\t\tOpcode Fx0A executed: Wait for a key press - (NOT PRESSED)\n")
		}
	}
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
// DT is set equal to the value of Vx.
func opc_chip8_FX15(x uint16) {
	DelayTimer = V[x]
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode Fx15 executed: Set delay timer = V[x(%d):%d]\n", x, V[x])
	}
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
// ST is set equal to the value of Vx.
func opc_chip8_FX18(x uint16) {
	SoundTimer = V[x]
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode Fx18 executed: Set sound timer = V[x(%d)]:%d\n",x, V[x])
	}
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
// The values of I and Vx are added, and the results are stored in I.
// ***
// Check FX1E (I = I + VX) buffer overflow. If buffer overflow, register
// VF must be set to 1, otherwise 0. As a result, register VF not set to 1.
// This undocumented feature of the Chip-8 and used by Spacefight 2091!
func opc_chip8_FX1E(x uint16) {
	if Debug {
		fmt.Printf("\t\tOpcode Fx1E executed: Add the value of V[x(%d)]:%d to I(%d)\n",x, V[x], I)
	}

	// *** Implement the undocumented feature used by Spacefight 2091
	if FX1E_spacefight2091 {
		if ( I + uint16(V[x]) > 0xFFF ) { //4095 - Buffer overflow
			V[0xF] = 1
			I = ( I + uint16(V[x]) ) - 4095
		} else {
			V[0xF] = 0
			I += uint16(V[x])
		}
	// Normal opcode pattern
	} else {
		I += uint16(V[x])
	}

	PC += 2
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
func opc_chip8_FX29(x uint16) {
	// Load CHIP-8 font. Start from Memory[0]
	I = uint16(V[x]) * 5
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode Fx29 executed: Set I(%X) = location of sprite for digit V[x(%d)]:%d (*5)\n", I, x, V[x])
	}
}

// Fx33 - LD B, Vx
// BCD - Binary Code hexadecimal
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
// set_BCD(Vx);
// Ex. V[x] = ff (maximum value) = 255
// memory[i+0] = 2
// memory[i+1] = 5
// memory[i+2] = 5
// % = modulus operator:
// 3 % 1 would equal zero (since 3 divides evenly by 1)
// 3 % 2 would equal 1 (since dividing 3 by 2 results in a remainder of 1).
func opc_chip8_FX33(x uint16) {
	Memory[I]   = V[x]  / 100
	Memory[I+1] = (V[x] / 10)  % 10
	Memory[I+2] = (V[x] % 100) % 10
	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode Fx33 executed: Store BCD representation of V[x(%d)]:%d in memory locations I(%X):%d, I+1(%X):%d, and I+2(%X):%d\n", x, V[x], I, Memory[I], I+1, Memory[I+1], I+2, Memory[I+2])
	}
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
//
// Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
// In the original CHIP-8 implementation, and also in CHIP-48, I is left incremented after this instruction had been executed. In SCHIP, I is left unmodified.
func opc_chip8_FX55(x uint16) {
	for i := uint16(0); i <= x; i++ {
		Memory[I+i] = V[i]
	}
	PC += 2

	// If needed, run the original Chip-8 opcode (not used in recent games)
	if Legacy_Fx55_Fx65 {
		I = I + x + 1
	}

	if Debug {
		fmt.Printf("\t\tOpcode Fx55 executed: Registers V[0] through V[x(%d)] in memory starting at location I(%d)\n",x, I)
	}
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
// The interpreter reads values from memory starting at location I into registers V0 through Vx.
//// I is set to I + X + 1 after operation²
//// ² Erik Bryntse’s S-CHIP documentation incorrectly implies this instruction does not modify
//// the I register. Certain S-CHIP-compatible emulators may implement this instruction in this manner.
func opc_chip8_FX65(x uint16) {
	for i := uint16(0); i <= x; i++ {
		V[i] = Memory[I+i]
	}

	PC += 2

	// If needed, run the original Chip-8 opcode (not used in recent games)
	if Legacy_Fx55_Fx65 {
		I = I + x + 1
	}

	if Debug {
		fmt.Printf("\t\tOpcode Fx65 executed: Read registers V[0] through V[x(%d)] from memory starting at location I(%X)\n",x, I)
	}
}

// ---------------------------- CHIP-8 undocumented instructions ---------------------------- //

// 02D8
// NON DOCUMENTED OPCODED, USED BY DEMO CLOCK Program
// LDA 02, I // Load from memory at address I into V[00] to V[02]
// Maybe an 0NNN 1802 instruction?
func opc_chip8_ND_02D8() {
	x := (Opcode & 0x0F00) >> 8

	if x != 2 {
		//Map if this opcode can receive a different value here
		fmt.Printf("\nProposital exit to map usage of 02D8 opcode\n")
		os.Exit(2)
	}

	V[0] = byte(I)
	V[1] = byte(I) + 1
	V[2] = byte(I) + 2

	PC += 2
	if Debug {
		fmt.Printf("\t\tOpcode 02DB executed (NON DOCUMENTED). - Load from memory at address I(%d) into V[0]= %d, V[1]= %d and V[2]= %d.\n", I, I , I+1, I+2)
	}
}