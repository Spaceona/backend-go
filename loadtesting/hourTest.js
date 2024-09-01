import { sleep, check } from 'k6';
import http, { StructuredRequestBody } from 'k6/http';

export let options = {
    vus: 105,
    duration: '5m',
    thresholds: {
        http_req_failed: ['rate<0.01'], // http errors should be less than 1%
        http_req_duration: ['p(95)<200'], // 95% of requests should be below 200ms
    },
};


export default () => {
    const response = http.get("http://localhost:3001/firmware/file/1-2-1", {
        headers: {
            "Authorization": "bearer eyJhbGciOiJIUzI1NiJ9.eyJkYXRhIjp7ImFwaVZlcnNpb24iOiJTcGFjZW9uYSAxLjIuMCIsImZpcm13YXJlVmVyc2lvbiI6IjEtMC0wIiwibWFjIjoiZTg6NjU6Mzg6NzU6MjI6OWQifSwiaWF0IjoxNzIzNzUwMTY2LCJleHAiOjE3MjM4MzY1NjZ9.mQUBC7CGC_r3w6-4HX230ogWCIM7gIq5Ez5nWu899Zc"
        }
    });
    check(response, {
        'status is 200': r => r.status === 200,
    });

    sleep(1);
};