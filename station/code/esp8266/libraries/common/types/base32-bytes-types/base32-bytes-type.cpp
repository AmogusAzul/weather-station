#include "base32-bytes-types.h"

Base32Bytes::Base32Bytes() :
Bytes::Bytes(BASE32_BYTES_LENGTH)
{
}

Base32Bytes::Base32Bytes(const uint8_t* data) :
Bytes::Bytes(BASE32_BYTES_LENGTH, data)
{
}