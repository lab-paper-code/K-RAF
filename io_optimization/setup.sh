#!/bin/bash
dataset_path="https://image-net.org/data/ILSVRC/2012/ILSVRC2012_img_val.tar"

mkdir images
cd images
wget $dataset_path
tar -xvf ILSVRC2012_img_val.tar
rm ILSVRC2012_img_val.tar

# 이미지를 카테고리 별로 분류
sh ../image_categorize.sh
cd ..

# venv 설정 및 requirements 설치
sudo apt install python3-venv
python3.8 -m venv env
. env/bin/activate
pip install --upgrade pip
pip install -r requirements.txt
