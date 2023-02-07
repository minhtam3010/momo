from flask import Flask, request, redirect
import random

app = Flask(__name__)

products = ["Portable Solar Charger", "Wireless Headphones", "Smart Watch", "Bluetooth Speaker", "Fitness Tracker"]

@app.route('/')
def index():
    selected_product = random.choice(products)
    return '''
        <h1>{0}</h1>
        <form action="/buy" method="post" target="_blank" enctype="application/x-www-form-urlencoded">
            <input type="hidden" name="product" value="{0}">
            <input type="submit" value="Buying via MoMo">
        </form>
    '''.format(selected_product)

@app.route('/buy', methods=['POST'])
def buy():
    import json
    import urllib.request
    import urllib
    import uuid
    import requests
    import hmac
    import hashlib
    product = request.form['product']
    # parameters send to MoMo get get payUrl
    endpoint = "https://test-payment.momo.vn/v2/gateway/api/create"
    partnerCode = "MOMO"
    accessKey = "F8BBA842ECF85"
    secretKey = "K951B6PE1waDMi640xX08PD3vg6EkVlz"
    orderInfo = "pay with MoMo"
    redirectUrl = "https://webhook.site/b3088a6a-2d17-4f8d-a383-71389a6c600b"
    ipnUrl = "https://webhook.site/b3088a6a-2d17-4f8d-a383-71389a6c600b"
    amount = "50000"
    orderId = str(uuid.uuid4())
    requestId = str(uuid.uuid4())
    requestType = "captureWallet"
    extraData = ""  # pass empty value or Encode base64 JsonString

    # before sign HMAC SHA256 with format: accessKey=$accessKey&amount=$amount&extraData=$extraData&ipnUrl=$ipnUrl
    # &orderId=$orderId&orderInfo=$orderInfo&partnerCode=$partnerCode&redirectUrl=$redirectUrl&requestId=$requestId
    # &requestType=$requestType
    rawSignature = "accessKey=" + accessKey + "&amount=" + amount + "&extraData=" + extraData + "&ipnUrl=" + ipnUrl + "&orderId=" + orderId + "&orderInfo=" + orderInfo + "&partnerCode=" + partnerCode + "&redirectUrl=" + redirectUrl + "&requestId=" + requestId + "&requestType=" + requestType

    # puts raw signature
    print("--------------------RAW SIGNATURE----------------")
    # signature
    h = hmac.new(bytes(secretKey, 'ascii'), bytes(rawSignature, 'ascii'), hashlib.sha256)
    signature = h.hexdigest()
    print("--------------------SIGNATURE----------------")

    # json object send to MoMo endpoint

    data = {
        'partnerCode': partnerCode,
        'partnerName': "Test",
        'storeId': "MomoTestStore",
        'requestId': requestId,
        'amount': amount,
        'orderId': orderId,
        'orderInfo': orderInfo,
        'redirectUrl': redirectUrl,
        'ipnUrl': ipnUrl,
        'lang': "vi",
        'extraData': extraData,
        'requestType': requestType,
        'signature': signature
    }
    print("--------------------JSON REQUEST----------------\n")
    data = json.dumps(data)

    clen = len(data)
    response = requests.post(endpoint, data=data, headers={'Content-Type': 'application/json', 'Content-Length': str(clen)})

    # f.close()
    print("--------------------JSON response----------------\n")

    return redirect(response.json()['payUrl'], code=302)

if __name__ == '__main__':
    app.run()
