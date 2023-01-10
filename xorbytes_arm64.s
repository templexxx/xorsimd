
// func bytesN(dst, a, b *byte, n int)
TEXT ·bytesN(SB), 4, $0
	MOVD  d+0(FP), R0
	MOVD  a+8(FP), R1
	MOVD  b+16(FP), R2
	MOVD  n+24(FP), R3
	TBNZ   $15, R3, not_aligned            // AND 15 & len, if not zero jump to not_aligned.

aligned:
	MOVD $0, R4 // position in slices

loop16b:
	VLD1 (R1)(R4*1), V0.16B   // XOR 16byte forwards.
	VLD1 (R2)(R4*1), V1.16B
	VEOR V0.16B, V1.16B, V0.16B
	VST1 V0.16B, (R0)(R4*1)
	ADD  $16, R4, R4
	CMP  R3, R4
	BNE  loop16b
	RET

loop_1b:
	SUB   $1, R3, R3           // XOR 1byte backwards.
	MOVB  (R1)(R3*1), DI
	MOVB  (R2)(R3*1), R4
	XORB  R4, DI
	MOVB  DI, (R0)(R3*1)
	TESTQ $7, R3           // AND 7 & len, if not zero jump to loop_1b.
	JNZ   loop_1b
	CMPQ  R3, $0           // if len is 0, ret.
	JE    ret
	TESTQ $15, R3          // AND 15 & len, if zero jump to aligned.
	JZ    aligned

not_aligned:
	TESTQ $7, R3           // AND $7 & len, if not zero jump to loop_1b.
	JNE   loop_1b
	SUBQ  $8, R3           // XOR 8bytes backwards.
	MOVD  (R1)(R3*1), DI
	MOVD  (R2)(R3*1), R4
	XORQ  R4, DI
	MOVD  DI, (R0)(R3*1)
	CMPQ  R3, $16          // if len is greater or equal 16 here, it must be aligned.
	JGE   aligned

ret:
	RET
	
// func bytes8(dst, a, b *byte)
TEXT ·bytes8(SB), 4, $0
	MOVD  d+0(FP), R0
	MOVD  a+8(FP), R1
	MOVD  b+16(FP), R2
	MOVD  (R1), R3
    MOVD  (R2), R4
    EOR   R3, R4
    MOVD  R4, (R0)
    RET

// func bytes16(dst, a, b *byte)
TEXT ·bytes16(SB), 4, $0
	MOVD  d+0(FP), R0
    MOVD  a+8(FP), R1
    MOVD  b+16(FP), R2
    VLD1  (R1), V0.B16
    VLD1  (R2), V1.B16
    VEOR  V0.B16, V1.B16, V0.B16
    VST1  V0.B16, (R0)
    RET
