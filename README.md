# TextOverImg
Naive implementation of an app that users can submit text and an image URL to.
The app returns the image with the text placed over it. Users can login to remove the image watermark.

The login is hard coded to username: test, password: test for now.

## How to build

```sh
cd frontend
npm run build
cd ..
go build .
```

## How to run
```sh
./TextOverImg
```

### Example queries to check endpoints:

curl -X POST -d "{\"url\": \"image-url-here.jpg\", \"text\": \"Inpsiration Quote Here!\"}" http://localhost:3000/image

curl -X POST -d "{\"username\": \"test\", \"password\": \"test\"}" http://localhost:3000/user

### TODO / improvements:
- implement a postgres db for users
    -hashed passwords using bcrypt
- better looking frontend
- store and serve images from a AWS S3 bucket / azure blob ?
