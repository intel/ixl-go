
TEXT ·ENQCMD(SB), $0-16
    MOVQ a+0(FP),AX
    MOVQ b+8(FP),BX
    SFENCE
    // ENQCMD AX,BX
    BYTE $0xf2;BYTE $0x0f;BYTE $0x38;BYTE $0xf8;BYTE $0x03;
    // CHECK ZF
    SETEQ AX
    MOVQ AX,c+16(FP)
	RET

TEXT ·UMONITOR(SB), $0-8
    MOVQ a+0(FP),AX
    BYTE $0xf3;
    BYTE $0x48;
    BYTE $0x0f;
    BYTE $0xae;
    BYTE $0xf0;

    RET
