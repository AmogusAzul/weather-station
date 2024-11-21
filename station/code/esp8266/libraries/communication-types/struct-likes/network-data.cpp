#include "struct-likes.h"
#include "../../common/types/base32-bytes-types/base32-bytes-types.h"
#include "../../common/types/string-bytes-types/string-bytes-types.h"  // For StringBytes type

#include <Arduino.h>

// Constructor with default values
NetworkData::NetworkData(const UInt32Bytes &stationId, const StringBytes &ip, const UInt32Bytes &port,
                         const StringBytes &ssid, const StringBytes &password)
    : stationId(), ip(), port(), ssid(), password() {

    Serial.print("Station ID: "); Serial.println(this->stationId.getNumber());
    Serial.print("SSID: "); Serial.println(this->ssid.getString().c_str());
    Serial.print("Password: "); Serial.println(this->password.getString().c_str());
    Serial.print("IP: "); Serial.println(this->ip.getString().c_str());
    Serial.print("Port: "); Serial.println(this->port.getNumber());
}

// Getter methods
UInt32Bytes NetworkData::getStationId() const {
    return stationId;
}

StringBytes NetworkData::getIp() const {
    return ip;
}

UInt32Bytes NetworkData::getPort() const {
    return port;
}

StringBytes NetworkData::getSsid() const {
    return ssid;
}

StringBytes NetworkData::getPassword() const {
    return password;
}