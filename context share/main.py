from flask import Flask, request
import requests
import os
import time
from tqdm import tqdm
import warnings
warnings.filterwarnings(action='ignore')

app = Flask(__name__)

img_path = "" # 추론할 이미지가 저장된 장소
storage_path = "" # 추론한 결과가 저장될 장소
pod_add = ""  # 실행을 공유할 가상 환경 IP


def requestToVirtualEnv(index, amount, limit):
    params = {
        'index': index,
        'amount': amount,
        'limit' : limit,
    }
    res = requests.post('http://' + pod_add, data=params)

    return res


def getFilesByIndex(start_index, amount):
    files = os.listdir(img_path)
    files = files[start_index:start_index + amount]
    return files

# from silence_tensorflow import silence_tensorflow
# silence_tensorflow()

from tensorflow.keras.applications import MobileNet
from keras.applications.inception_v3 import preprocess_input
from tensorflow.keras.utils import img_to_array
from tensorflow.keras.utils import load_img
from keras.applications import imagenet_utils
import numpy as np

def predict_image(filename, file, cnt):
    model = MobileNet(weights="imagenet")
    inputShape = (224, 224)
    path = img_path + filename
    img = load_img(path, target_size=inputShape)
    image = img_to_array(img)
    image = np.expand_dims(image, axis=0)
    image = preprocess_input(image)

    predicts = model.predict(image)
    P = imagenet_utils.decode_predictions(predicts)

    for (i, (imagenetID, label, prob)) in enumerate(P[0]):
        file.write("{}. {}: {:.2f}% \n".format(cnt, label, prob * 100))
        # print("{}: {:.2f}%".format(label, prob * 100))
        return imagenetID, label



@app.route('/share_pi', methods=['GET', 'POST'])
def contextSharePI():
    try:
        index = int(request.form.get('index')) # 시작 인덱스 처음 시작시 0
        amount = int(request.form.get('amount'))# ex) 1000장씩 context share
        limit = int(request.form.get('limit')) # ex) 추론할 이미지가 전체 10000장일 경우 10000
        print(f'index : {index}, amount : {amount}, limit:{limit}')
        start_index = index * amount

        mode = "a"
        cnt = 1

        if start_index < limit:
            files = getFilesByIndex(start_index, amount)
            result_file_path = storage_path + "result.txt"

            if index == 0:
                if os.path.isfile(result_file_path):
                    os.unlink(result_file_path)
                    mode = "w"
            else:
                mode = "a"
            f = open(result_file_path, mode)
            # 3. 이미지 추론
            # 4. 추론 결과는 공유된 가상 저장소에 저장하기

            start = time.time()
            for file in tqdm(files):
                if file.endswith('.JPEG') or file.endswith('.JPG') or file.endswith('.PNG'):
                    predict_image(file, f, start_index + cnt)
                    cnt += 1
            f.close()
            print("Elapsed Time : %s" % (time.time() - start))
            # 5. 1000장 추론이 끝나면 그 이후의 context를 pod으로 넘기기
            index += 1

            res = requestToVirtualEnv(index, amount, limit)

            if res.ok:
                print(res.text)
            else:
                print("Failed")
                print(res.text)
        else:
            print("")
        return "Success"
    except Exception as e:
        return str(e)




if __name__ == '__main__':
    app.run(host="0.0.0.0", port=60000, debug=True)


