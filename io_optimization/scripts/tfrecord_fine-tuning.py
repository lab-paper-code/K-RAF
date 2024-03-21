import sys
import re

import tensorflow as tf
import glob

def get_args(argv: list) -> dict:
    arg_parser = re.compile(r'--(?P<option_name>[a-zA-Z_]+)=(?P<value>[\S]+)')
    valid_options = [
        'dataset_path',
        'dataset_size',
        'checkpoint_path',
        'log_path',
        'batch_size',
        'epochs',
        'base_model'
    ]
    settings = {
        # Default setting
        'dataset_path': '.',
        'dataset_size': 50000,
        'batch_size': 32,
        'epochs': 1,
        'base_model': None
    }

    for arg in argv:
        parsing_result = arg_parser.search(arg)
        if parsing_result is None:
            print('[ERR] Wrong option:', arg, '(Ignored)')
        else:
            option_name = parsing_result.group('option_name')
            option_value = parsing_result.group('value')
            if option_name not in valid_options:
                print('[ERR] Invalid option:', option_name, '(Ignored)')
            else:
                if option_name in ['dataset_size', 'batch_size', 'epochs']:
                    option_value = int(option_value)
                settings[option_name] = option_value
    
    return settings

def set_callbacks(settings:dict) -> list:
    callbacks = []
    
    if 'checkpoint_path' in settings.keys():
        if not settings['checkpoint_path'].endswith('/'): settings['checkpoint_path'] += '/'
        checkpoint_filename = 'checkpoint-epoch-{}.h5'.format(settings['epochs'])
        checkpoint = ModelCheckpoint(
            checkpoint_filename,
            monitor='accuracy',
            verbose=1,
            save_best_only=True,
            mode='auto'
            )
        callbacks.append(checkpoint)
    
    if 'log_path' in settings.keys():
        total_steps = settings['dataset_size']//settings['batch_size']
        # 중간의 100 Steps를 프로파일링
        mid = total_steps // 2
        start_step = mid - 49
        end_step = mid + 50
        if start_step <= 0:
            start_step = 1
        if end_step > total_steps:
            end_step = total_steps

        profile = tf.keras.callbacks.TensorBoard(
            log_dir=settings['log_path'],
            histogram_freq=1,
            profile_batch=[start_step, end_step]
        )
        callbacks.append(profile)
    
    return callbacks

# TFRecord 파일 파싱 함수
def parse_tfrecord(serialized_example):
    feature_description = {
        'name': tf.io.FixedLenFeature([], tf.string),
        'category': tf.io.FixedLenFeature([], tf.int64),
        'data': tf.io.FixedLenFeature([], tf.string)
    }
    example = tf.io.parse_single_example(serialized_example, feature_description)
    
    image = tf.io.decode_jpeg(example['data'], channels=3)
    image = tf.image.resize(image, [224, 224])
    image = (image / 127.5) - 1
    label = tf.one_hot(example['category'], depth=1000)
    
    return image, label

def main():
    settings = get_args(sys.argv[1:])
    callbacks = set_callbacks(settings)

    tfrecord_list = glob.glob(settings['dataset_path'] + '/*.tfrecord')

    dataset = tf.data.TFRecordDataset(tfrecord_list, num_parallel_reads=None)
    dataset = dataset.map(parse_tfrecord)
    dataset = dataset.shuffle(seed=42, buffer_size=256)
    dataset = dataset.batch(settings['batch_size'])

    base_model = tf.keras.applications.MobileNet(weights=settings['base_model'], include_top=False, input_shape=(224, 224, 3))
    # base_model.trainable = False
    x = base_model.output
    x = tf.keras.layers.GlobalAveragePooling2D()(x)
    x = tf.keras.layers.Dense(1000, activation='softmax')(x)
    model = tf.keras.Model(inputs=base_model.input, outputs=x)

    model.compile(optimizer=tf.keras.optimizers.Adam(learning_rate=0.0001), loss='categorical_crossentropy', metrics=['accuracy'])

    model.fit(dataset, epochs=settings['epochs'], callbacks=callbacks)

main()
