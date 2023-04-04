const HOST = 'copilot-proxy.githubusercontent.com';

export default {
  async fetch(request, env) {
    try {
      const url = new URL(request.url);
      url.hostname = HOST;

      const forwardedRequest = new Request(url, {
        method: request.method,
        headers: request.headers,
        body: request.body
      });

      const response = await fetch(forwardedRequest);
      return response;
    } catch(e) {
      return new Response(err.stack, { status: 500 })
    }
  }
}