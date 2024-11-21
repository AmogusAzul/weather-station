#include "communication-bytes-types.h"
#include "../../safety/safety.h"



TokenBytes::TokenBytes(const uint8_t* data) :
Bytes::Bytes(safety::TOKEN_LENGTH, data)
{
}

TokenBytes::TokenBytes(const std::string str) :
Bytes::Bytes(safety::TOKEN_LENGTH)
{
    setByteString(str);
}

