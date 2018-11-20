# publishapk
Lightweight utility for publishing APKs to the PlayStore

## Usage

Usage of ./publishapk:  
  -apk="": path to the apk which should be uploaded.  
  -debug=false: flag indicating whether debug output should be printed.  
  -email="": email address of the service account.  
  -key="": path to the key.pem file for the service account OR base64 encoded private key.  
  -pkg="": package name of the application.  
  -track="": track to which the apk should be pushed (e.g.: alpha, beta, internal, production)  
  
All of those parameter can also be specified as environment variable, e.g. `APK=/tmp/some.apk ./publishapk`
