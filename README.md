# TextOverImg
Naive implementation of an app that users can submit text and an image URL to.
The app returns the image with the text placed over it.


URL request example:
curl -X POST -d "{\"url\": \"https://upload.wikimedia.org/wikipedia/commons/3/3d/Forstarbeiten_in_%C3%96sterreich.JPG\"}" http://localhost:3000/image

TODO:
- Add text field to the POST request
- Add error handling	
	- invalid POST request: url or text missing
	- unable to get image from url
	- text too long?
- Add watermark / premium, user db and login/logout functions
- Add a frontend in Vue.js
- (Extra QoL) Improve the save file / load file approach, directly get an image.Image from the response body
