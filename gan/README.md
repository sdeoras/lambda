# MNIST GAN
The code here deploys a Google cloud function to show a gallery
of 10 randomly generated MNIST GAN images.

## Arch
* Model was trained in Python and exported as a protocol buffer graph along with frozen weights
* Such model was then imported in Go and put together in a Google Cloud function with Go runtime
* Frontend was designed in Adobe Lightroom
* API was designed in Google protocol buffer (proto3)

## Screenshot
![mnist-gan](https://github.com/sdeoras/lambda/raw/master/gan/art/mnist-gan-gallery.png "MNIST GAN Gallery")
