import { sleep, check } from 'k6';
import http, { StructuredRequestBody } from 'k6/http';
import exec from 'k6/execution';
const max_machines = 1000
export let options = {
    stages: [
        // { duration: '5m', target: 100 }, // traffic ramp-up from 1 to 100 users over 5 minutes.
        // { duration: '5m', target: 100 }, // stay at 100 users for 30 minutes
        // { duration: '5m', target: 1000 }, // stay at 100 users for 30 minutes
        // { duration: '5m', target: 1000 }, // stay at 100 users for 30 minutes
        { duration: '10m', target: max_machines }, // stay at 100 users for 30 minutes
        // { duration: '5m', target: 1000 }, // stay at 100 users for 30 minutes
        // { duration: '5m', target: 100 }, // stay at 100 users for 30 minutes
        // { duration: '5m', target: 0 }, // ramp-down to 0 users
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'], // http errors should be less than 1%
        http_req_duration: ['p(95)<200'], // 95% of requests should be below 200ms
    },
};
const authToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImVtYWlsIjp7ImlkIjoiMTAyMTc5ODE1NTQ3MzkyMjE2NzgwIiwiZW1haWwiOiJuaWNrLmxlc2xpZTMwM0BnbWFpbC5jb20iLCJ2ZXJpZmllZF9lbWFpbCI6dHJ1ZSwibmFtZSI6Ik5pY2sgTGVzbGllIiwiZ2l2ZW5fbmFtZSI6Ik5pY2siLCJmYW1pbHlfbmFtZSI6Ikxlc2xpZSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NJRVJzZkxpYnRKN1FhcHRTWlgxMUJwaXpQaWJTckpZMXF3cm9LR09xWnEwNFNPUkpRVj1zOTYtYyJ9LCJHb29nbGVUb2tlbiI6eyJhY2Nlc3NfdG9rZW4iOiJ5YTI5LmEwQWNNNjEyeWJCanFQR2xaUk9FYS0wN1ZUdk9rOWJ4THdzT3R6VmFXcTMxVHpPZ0llTFBHSjM4M3ZTTkxHSWtMT2h5MlZkYk41RVZIeFRFQmR5b014ZHpjQ24yWU1WZlJMd2xrUkVVNDN5UHBRb1QxUDNVN1hQaFF6THpoaGRWdGcxV2p2MmlDN1dHU1dibUpzNmRKRDR6NndxMWlNd2hTU21DQWVOT001YUNnWUtBWmNTQVJFU0ZRSEdYMk1pLVRnbnpxS0dvTWpoV1VFZS1HMUE0dzAxNzUiLCJleHBpcmVzX2luIjozNTk5LCJyZWZyZXNoX3Rva2VuIjoiMS8vMDFJT3B4Y3dkRUk2R0NnWUlBUkFBR0FFU053Ri1MOUlybFdrbElJUEFiNm83X1VDNVNPY1BnbXE1emE4Zk00dG9MQlY2WldCWmVMMllTcHdqY0JxWDZrLWZWWWxwMllyYjk1dyIsInNjb3BlIjoib3BlbmlkIGh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL2F1dGgvdXNlcmluZm8ucHJvZmlsZSBodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9hdXRoL3VzZXJpbmZvLmVtYWlsIiwidG9rZW5fdHlwZSI6IkJlYXJlciJ9fSwiZXhwIjoxNzI2OTQ1MDUwfQ.sa4QIgTyxDhC8LDAtNVCaCPchnPIXq5BH6oth-b6joM"

export async function setup() {
    const machines = [];
    for (let i = 0; i < max_machines/2; i++) {
        machines.push({"number":i,"type":"Washer"})
    }
    for (let i = 0; i < max_machines/2; i++) {
        machines.push({"number":max_machines/2+i,"type":"Dryer"})
    }
    
    // 2. setup code
    const clientPayload = {
        "client_name":"nicks laundry",
        "buildings":[
            {
                "building_name":"test laundry",
                "machines": machines
            }
        ]
    }
    const res =  http.post("http://localhost:3000/onboard/client",JSON.stringify(clientPayload),{
        headers: {
            "Authorization": "bearer " + authToken
        }
    });
    console.log(res.status)
    const json = await res.json()
    console.log(json)
    const key = json.key;
    if(key === undefined) {
        return;
    }
    const boardMappings = []
    for (let i = 0; i < max_machines; i++) {
        const macAddress = "loadTest" + i.toString();
        const boardBody = JSON.stringify({
            mac_address:macAddress,
            client_name:"nicks laundry",
            client_key:key
        })
        const boardRes = http.post("http://localhost:3000/onboard/board",boardBody,{
            headers: {
                "Authorization": "bearer " + authToken
            }
        });
        const boardJson = await boardRes.json()
        boardMappings.push({mac_address:macAddress,machine_id:i})
        if(i % 20 == 0) {
            sleep(0.2);
        }
    }
    const assignRes = http.post("http://localhost:3000/admin/machine/assign",JSON.stringify({mappings:boardMappings}),{
        headers: {
            "Authorization": "bearer " + authToken
        }
    })
    const assignJson = await assignRes.json()
    console.log(assignJson)
}


export default () => {
    const payload = JSON.stringify({
        "mac_address":"loadTest"+exec.vu.idInTest,
        "firmware_version":"1-2-1",
        "status":Math.random() < 0.5,
        "StatusChanged":true,
        "timeBetweenChange":1,
        "confidence":100
    });

    const response = http.post("http://localhost:3000/status/update",payload  ,{
        headers: {
            "Authorization": "bearer " + authToken
        }
    });
    check(response, {
        'status is 200': r => r.status === 200,
    });

    sleep(Math.random() * 5);
};