package main

import (
   "fmt"
   "github.com/elt/googleplay"
   "github.com/elt/rosso/http"
   "io"
   "os"
   "time"
)

func (f flags) do_auth(dir string) error {
   res, err := googleplay.New_Auth(f.email, f.password)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return res.Create(dir + "/auth.txt")
}

func do_device(dir, platform string) error {
   res, err := googleplay.Phone.Checkin(platform)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   fmt.Printf("Sleeping %v for server to process\n", googleplay.Sleep)
   time.Sleep(googleplay.Sleep)
   return res.Create(dir + "/" + platform + ".bin")
}

func (f flags) do_delivery(head *googleplay.Header) error {
   download := func(ref, name string) error {
      res, err := googleplay.Client.Redirect(nil).Get(ref)
      if err != nil {
         return err
      }
      defer res.Body.Close()
      file, err := os.Create(name)
      if err != nil {
         return err
      }
      defer file.Close()
      pro := http.Progress_Bytes(file, res.ContentLength)
      if _, err := io.Copy(pro, res.Body); err != nil {
         return err
      }
      return nil
   }
   del, err := head.Delivery(f.app, f.version)
   if err != nil {
      return err
   }
   file := googleplay.File{f.app, f.version}
   for _, split := range del.Split_Data() {
      ref, err := split.Download_URL()
      if err != nil {
         return err
      }
      id, err := split.ID()
      if err != nil {
         return err
      }
      if err := download(ref, file.APK(id)); err != nil {
         return err
      }
   }
   for _, add := range del.Additional_File() {
      ref, err := add.Download_URL()
      if err != nil {
         return err
      }
      typ, err := add.File_Type()
      if err != nil {
         return err
      }
      if err := download(ref, file.OBB(typ)); err != nil {
         return err
      }
   }
   ref, err := del.Download_URL()
   if err != nil {
      return err
   }
   return download(ref, file.APK(""))
}

func (f flags) do_header(dir, platform string) (*googleplay.Header, error) {
   var head googleplay.Header
   err := head.Open_Auth(dir + "/auth.txt")
   if err != nil {
      return nil, err
   }
   if err := head.Auth.Exchange(); err != nil {
      return nil, err
   }
   if err := head.Open_Device(dir + "/" + platform + ".bin"); err != nil {
      return nil, err
   }
   head.Single = f.single
   return &head, nil
}

func (f flags) do_details(head *googleplay.Header) ([]byte, error) {
   detail, err := head.Details(f.app)
   if err != nil {
      return nil, err
   }
   return detail.MarshalText()
}
