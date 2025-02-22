// server.ts
import { serveFile } from "https://deno.land/std@0.177.0/http/file_server.ts";

const handler = async (request: Request): Promise<Response> => {
  const url = new URL(request.url);
  if (url.pathname === "/") {
    return await serveFile(request, "./index.html");
  } else {
    return new Response("404 Not Found", { status: 404 });
  }
};

console.log("HTTP webserver running. Access it at: http://localhost:8000/");
Deno.serve(handler);
