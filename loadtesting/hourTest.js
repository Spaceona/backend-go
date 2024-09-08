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

export async function setup() {
    const authRes =  http.post("http://localhost:3001/auth/user",JSON.stringify({}));
    const authJson = await authRes.json()
    const token = authJson.Token;

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
    const res =  http.post("http://localhost:3001/onboard/client",JSON.stringify(clientPayload),{
        headers: {
            "Authorization": "bearer " + token
        }
    });
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
        const boardRes = http.post("http://localhost:3001/onboard/board",boardBody,{
            headers: {
                "Authorization": "bearer " + token
            }
        });
        const boardJson = await boardRes.json()
        boardMappings.push({mac_address:macAddress,machine_id:i})
        if(i % 20 == 0) {
            sleep(0.2);
        }
    }
    const assignRes = http.post("http://localhost:3001/admin/machine/assign",JSON.stringify({mappings:boardMappings}),{
        headers: {
            "Authorization": "bearer " + token
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

    const response = http.post("http://localhost:3001/status/update",payload  ,{
        headers: {
            "Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7Im1hY19hZGRyZXNzIjoiNjQ6RTg6MzM6ODY6REI6QTQiLCJmaXJtd2FyZV92ZXJzaW9uIjoiMS0wLTAifSwiZXhwIjoxNzI1ODUwMzc0fQ.TxjD1KretiqsiVvpffooPicJkiUR5WV6YvvCb5WZsI8"
        }
    });
    check(response, {
        'status is 200': r => r.status === 200,
    });

    sleep(Math.random() * 2);
};