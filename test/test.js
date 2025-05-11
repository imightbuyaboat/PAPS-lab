import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10,
  duration: '1m',
  thresholds: {
    http_req_duration: ['p(95)<600'], 
  },
};

export function setup() {
  return { sessionId: '6d0768e2-e2ed-415b-b72d-1501549a2c73' }; 
}

export default function (data) {
  const headers = {
    'Cookie': `session_id=${data.sessionId}`,
    'Content-Type': 'application/x-www-form-urlencoded',
  };

  let res = http.get('https://sweet-things-melt.loca.lt/login', null);
  check(res, { 'GET /login': (r) => r.status === 200 });

  res = http.get('https://sweet-things-melt.loca.lt/register', null);
  check(res, { 'GET /register': (r) => r.status === 200 });

  res = http.get('https://sweet-things-melt.loca.lt/', { headers });
  check(res, { 'GET /': (r) => r.status === 200 });

  res = http.post('https://sweet-things-melt.loca.lt/search', 'organization=&city=&phone=', { headers });
  check(res, { 'POST /search': (r) => r.status === 200 });

  sleep(1);
}