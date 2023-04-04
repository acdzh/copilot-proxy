const http = require('http');
const https = require('https');
const url = require('url');

const clashUrl = 'http://127.0.0.1:7890';

const proxy = http.createServer((req, res) => {
  const start = new Date();

  // 如果需要代理可取消注释
  // const clashAgent = new https.Agent({
  //   proxy: url.parse(clashUrl),
  //   keepAlive: true,
  //   keepAliveMsecs: 60000,
  //   maxSockets: Infinity,
  //   timeout: 0,
  // });

  const options = {
    protocol: 'https:',
    host: 'copilot-proxy.githubusercontent.com',
    port: 443,
    path: req.url,
    method: req.method,
    headers: req.headers,
    // agent: clashAgent,
  };

  const proxyReq = https.request(options, (proxyRes) => {
    res.writeHead(proxyRes.statusCode, proxyRes.headers);
    proxyRes.pipe(res);
    proxyRes.on('end', () => {
      const elapsed = new Date() - start;
      console.log(`[${req.method}] ${req.url} Completed in ${elapsed}ms`);
    });
  });

  proxyReq.on('error', (err) => {
    console.error(err);
    res.writeHead(500);
    res.end('Internal Server Error');
  });

  req.pipe(proxyReq);
});

proxy.listen(9394);
console.log('Server started on http://127.0.0.1:9394');
