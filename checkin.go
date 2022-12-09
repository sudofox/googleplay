package googleplay

import (
   "bytes"
   "github.com/elt/rosso/protobuf"
   "net/http"
   "strconv"
   "time"
)

const Sleep = 4 * time.Second

// These can use default values, but they must all be included
type Config struct {
   GL_ES_Version uint64
   GL_Extension []string
   Has_Five_Way_Navigation uint64
   Has_Hard_Keyboard uint64
   Keyboard uint64
   Navigation uint64
   New_System_Available_Feature []string
   Screen_Density uint64
   Screen_Layout uint64
   System_Shared_Library []string
   Touch_Screen uint64
}

var Phone = Config{
   New_System_Available_Feature: []string{
      // app.source.getcontact
      "android.hardware.location.gps",
      // br.com.rodrigokolb.realdrum
      "android.software.midi",
      // com.app.xt
      "android.hardware.camera.front",
      // com.clearchannel.iheartradio.controller
      "android.hardware.microphone",
      // com.google.android.apps.walletnfcrel
      "android.software.device_admin",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      "android.hardware.wifi",
      // com.illumix.fnafar
      "android.hardware.sensor.gyroscope",
      // com.madhead.tos.zh
      "android.hardware.sensor.accelerometer",
      // com.miHoYo.GenshinImpact
      "android.hardware.opengles.aep",
      // com.pinterest
      "android.hardware.camera",
      "android.hardware.location",
      "android.hardware.screen.portrait",
      // com.supercell.brawlstars
      "android.hardware.touchscreen.multitouch",
      // com.sygic.aura
      "android.hardware.location.network",
      // com.xiaomi.smarthome
      "android.hardware.bluetooth",
      "android.hardware.bluetooth_le",
      "android.hardware.camera.autofocus",
      "android.hardware.usb.host",
      // kr.sira.metal
      "android.hardware.sensor.compass",
      // org.thoughtcrime.securesms
      "android.hardware.telephony",
      // org.videolan.vlc
      "android.hardware.screen.landscape",
   },
   System_Shared_Library: []string{
      // com.amctve.amcfullepisodes
      "org.apache.http.legacy",
      // com.binance.dev
      "android.test.runner",
      // com.miui.weather2
      "global-miui11-empty.jar",
   },
   GL_Extension: []string{
      // com.instagram.android
      "GL_OES_compressed_ETC1_RGB8_texture",
      // com.kakaogames.twodin
      "GL_KHR_texture_compression_astc_ldr",
   },
   // com.axis.drawingdesk.v3
   // valid range 0x3_0001 - 0x7FFF_FFFF
   GL_ES_Version: 0xF_FFFF,
}

// A Sleep is needed after this.
func (c Config) Checkin(native_platform string) (*Response, error) {
   req_body := protobuf.Message{
      4: protobuf.Message{ // checkin
         1: protobuf.Message{ // build
            // sdkVersion
            // multiple APK valid range 14 - 0x7FFF_FFFF
            // single APK valid range 14 - 28
            10: protobuf.Varint(28),
         },
         18: protobuf.Varint(1), // voiceCapable
      },
      // version
      // valid range 2 - 3
      14: protobuf.Varint(3),
      18: protobuf.Message{ // deviceConfiguration
         1: protobuf.Varint(c.Touch_Screen),
         2: protobuf.Varint(c.Keyboard),
         3: protobuf.Varint(c.Navigation),
         4: protobuf.Varint(c.Screen_Layout),
         5: protobuf.Varint(c.Has_Hard_Keyboard),
         6: protobuf.Varint(c.Has_Five_Way_Navigation),
         7: protobuf.Varint(c.Screen_Density),
         8: protobuf.Varint(c.GL_ES_Version),
         11: protobuf.String(native_platform),
      },
   }
   for _, library := range c.System_Shared_Library {
      // .deviceConfiguration.systemSharedLibrary
      req_body.Get(18).Add_String(9, library)
   }
   for _, extension := range c.GL_Extension {
      // .deviceConfiguration.glExtension
      req_body.Get(18).Add_String(15, extension)
   }
   for _, name := range c.New_System_Available_Feature {
      // .deviceConfiguration.newSystemAvailableFeature
      req_body.Get(18).Add(26, protobuf.Message{
         1: protobuf.String(name),
      })
   }
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/checkin",
      bytes.NewReader(req_body.Marshal()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-protobuffer")
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   return &Response{res}, nil
}

type Device struct {
   protobuf.Message
}

func (d Device) ID() (uint64, error) {
   return d.Get_Fixed64(7)
}

type Native_Platform map[int64]string

var Platforms = Native_Platform{
   // com.google.android.youtube
   0: "x86",
   // com.miui.weather2
   1: "armeabi-v7a",
   // com.kakaogames.twodin
   2: "arm64-v8a",
}

func (n Native_Platform) String() string {
   b := []byte("nativePlatform")
   for key, val := range n {
      b = append(b, '\n')
      b = strconv.AppendInt(b, key, 10)
      b = append(b, ": "...)
      b = append(b, val...)
   }
   return string(b)
}
