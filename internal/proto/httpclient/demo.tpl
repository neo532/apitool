message {{ .ReplyName }} { 
    int32 code = 1;
    string message = 2;
    {{ .ReplyType }} data = 3;
}
