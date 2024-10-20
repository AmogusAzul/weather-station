#ifndef STRUCT_LIKES_H
#define STRUCT_LIKES_H

#include "../../../common/types/base32-bytes-types/base32-bytes-types.h"
#include "../../../common/types/string-bytes-types/string-bytes-types.h"

class NetworkData {
public:
    // Constructor with default values
    NetworkData(const UInt32Bytes &stationId = UInt32Bytes(0), 
                const StringBytes &ip = StringBytes("0.0.0.0"),
                const StringBytes &port = StringBytes("8080"),
                const StringBytes &ssid = StringBytes("defaultSSID"),
                const StringBytes &password = StringBytes("defaultPassword"));

    // Getter methods
    UInt32Bytes getStationId() const;
    StringBytes getIp() const;
    StringBytes getPort() const;
    StringBytes getSsid() const;
    StringBytes getPassword() const;

private:
    UInt32Bytes stationId;  // Store station ID as UInt32Bytes
    StringBytes ip;
    StringBytes port;
    StringBytes ssid;
    StringBytes password;
};

#endif  // STRUCT_LIKES_H