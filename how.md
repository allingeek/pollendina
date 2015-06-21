# Generate root certificate and run CA
docker run -d --name pollendina_ca -p 33004:33004 -v "$PWD":/opt/pollendina/ pollendina/pollendina

# Copy the certificate to example and create pollendina client image
copy 

build


run 
