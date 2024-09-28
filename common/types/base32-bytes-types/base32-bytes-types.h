#ifndef BASE32_BYTES_TYPES_H
#define BASE32_BYTES_TYPES_H

#include "../bytes-types.h"

const size_t BASE32_BYTES_LENGTH = 4;

class Base32Bytes : public Bytes {
public:

    // Constructors w/o a base byte array
    Base32Bytes(); 
    Base32Bytes(const uint8_t* data); 

};

#endif // BASE32BYTES_TYPES_H