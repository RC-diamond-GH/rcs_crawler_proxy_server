import requests

proxies = {
    'http': 'http://localhost:8080'
}

url = 'http://www.rcdiamondgh.cc/db.html?key=abc'

header = {
    "no-cache": "true"
}

if __name__ == '__main__':
    req = requests.get(url=url, proxies=proxies, headers=header, data="Long may the sunshine upon this lord of cinder!")

    print(req.status_code)
    for k, v in req.headers.items():
        print(f'{k}: {v}')
    
    if len(req.content) > 1000:
        with open('output.bin', 'wb') as j:
            print('has received {} bytes'.format(j.write(req.content)))
    else:
        print(req.content)
