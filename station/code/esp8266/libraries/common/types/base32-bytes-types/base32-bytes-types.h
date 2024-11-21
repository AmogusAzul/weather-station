#ifndef BASE32_BYTES_TYPES_H
#define BASE32_BYTES_TYPES_H

#include "../bytes-types.h"

constexpr uint8_t BASE32_BYTES_LENGTH = 4;

class Base32Bytes : public Bytes {
public:

    // Constructors w/o a base byte array
    Base32Bytes(); 
    Base32Bytes(const uint8_t* data);


};

class Int32Bytes : public Base32Bytes {
public:
    // Constructors for getting data from a number or an array
    Int32Bytes(const int32_t number);
    Int32Bytes(const uint8_t* data);

    

    // I/O functions for numbers
    void setNumber(const int32_t number);
    int32_t getNumber() const;

    static const int ID;

};

class UInt32Bytes : public Base32Bytes {
public:
    // Constructors for getting data from a number or an array
    UInt32Bytes(const uint32_t number=0);
    UInt32Bytes(const uint8_t* data);

    

    // I/O functions for numbers
    void setNumber(const uint32_t number);
    uint32_t getNumber() const;
};

class Float32Bytes : public Base32Bytes {
public:

    // Constructors for getting data from a number or an array
    Float32Bytes(const float number);
    Float32Bytes(const uint8_t* data);
    

    // I/O functions for numbers
    void setNumber(const float number);
    float getNumber() const;

};

#endif // BASE32BYTES_TYPES_H