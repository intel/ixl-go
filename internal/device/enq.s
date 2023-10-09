TEXT ·endcmdWithRetry(SB), $0-16
    MOVQ a+0(FP),AX
    MOVQ b+8(FP),BX
    SFENCE
loop:
    // ENQCMD EAX,[EBX]
    BYTE $0xf2;BYTE $0x0f;BYTE $0x38;BYTE $0xf8;BYTE $0x03;
    // CHECK ZF
    SETEQ CX
    CMPB CX,$0x00
    JE finish
    PAUSE
    JMP loop
finish:
    MOVQ CX,c+16(FP)
    RET

TEXT ·enqcmd(SB), $0-16
    MOVQ a+0(FP),AX
    MOVQ b+8(FP),BX
    SFENCE
    // ENQCMD AX,BX
    BYTE $0xf2;BYTE $0x0f;BYTE $0x38;BYTE $0xf8;BYTE $0x03;
    // CHECK ZF
    SETEQ AX
    MOVQ AX,c+16(FP)
    RET

TEXT ·movdir64b(SB), $0-16
    MOVQ a+0(FP),AX
    MOVQ b+8(FP),DX
    SFENCE
    // movdir64b rax,[rdx]
    BYTE $0x66
    BYTE $0x0f
    BYTE $0x38
    BYTE $0xf8
    BYTE $0x02
    RET 


TEXT ·waitForComplete(SB), $0-16
    MOVQ a+0(FP),AX
loop:
    MOVB 0(AX), BX
    LFENCE
    CMPB BX, $0x00
    JNE finish
    PAUSE
    JMP loop
finish:
    ANDB $0x1f,BX
    MOVQ BX,b+8(FP)
	RET

; TEXT ·UMONITOR_UMWAIT(SB), $0-16
; loop:
;     MOVQ a+0(FP),AX
;     // UMONITOR EAX
;     BYTE $0xf3;BYTE $0x0f;BYTE $0xae; BYTE $0xf0;

;     // RDTSC get current timestamp
;     // result in EDX:EAX
;     RDTSC

;     // move timestamp to RCX
;     MOVL DX,CX
;     SHLQ $32,CX
;     ADDL AX,CX
;     // use current timestamp + 200 as umwait deadline 
;     ADDQ $200,CX

;     // UMWAIT use EDX:EAX as deadline
;     MOVL CX,AX
;     SHRQ $32,CX
;     MOVL CX,DX

;     // set UMWAIT control bit
;     MOVQ $1,CX

;     // UMWAIT ECX
;     BYTE $0xf2; BYTE $0x0f; BYTE $0xae; BYTE $0xf1	

;     // check whether *addr has changed
;     MOVQ a+0(FP),AX
;     MOVB 0(AX), BX
;     LFENCE
;     CMPB BX, $0x00
;     JNE finish
;     JMP loop

; finish:
;     ANDB $0x1f,BX
;     MOVQ BX,b+8(FP)
;     RET
