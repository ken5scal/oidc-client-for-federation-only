const http = require('http');
const { exit } = require('process');
const port = (process.env.PORT || 8000 )
const hostname = (process.env.HOSTNAME || 'localhost')
const tenantID = process.env.AZURE_AD_TENANT_ID
const clientID = process.env.AZURE_AD_CLIENT_ID

if (tenantID == '' || clientID == '') {
    exit(1)
}

const app = http.createServer((req, res) => {
    var path = req.url.substring(1)
    if (path == 'redirect') {
        let body = '';
        req.on('data', (chunk) => {
            body += chunk.toString();
          }).on('end', () => {
            const firstParam = body.split('&')[0];
            const secondParam = body.split('&')[1];
            switch (firstParam.split("=")[0]) {
                case "id_token":
                    res.writeHead(200).end(firstParam.split("=")[1]);
                    break;
                case "error":
                    res.writeHead(401).end(secondParam.split("=")[1])
                    break;
            }
          });
    } else if (path == '') {
        console.log('redirect')
        res.writeHead(301, {
            location: `https://login.microsoftonline.com/${tenantID}/oauth2/v2.0/authorize?client_id=${clientID}&response_type=id_token&redirect_uri=http%3A%2F%2F${hostname}%3A${port}%2Fredirect&response_mode=form_post&scope=openid&state=123&nonce=123`
        })
        res.end();
    }
});

app.listen(port, hostname, () => {});