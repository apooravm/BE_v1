import requests
import json

url = "https://multi-serve.onrender.com/api/journal/"

payload = json.dumps({
  "username": "mrepig",
  "password": "XXXX",
  "log": "nother one 😋",
  "title": "random note",
  "tags": [
    "random thoughts",
    "todo"
  ]
})
headers = {
  'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)
