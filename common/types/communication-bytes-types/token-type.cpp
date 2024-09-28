#include "communication-bytes-types.h"
#include "../../safety/safety.h"



TokenBytes::TokenBytes(uint8_t* data) :
Bytes::Bytes(TOKEN_LENGTH, data)
{
}

TokenBytes::TokenBytes(std::string str) :
Bytes::Bytes(TOKEN_LENGTH)
{
    setByteString(str);
}

