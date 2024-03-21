import tensorflow as tf
import os
import sys
import re

def get_args(argv: list) -> dict:
    arg_parser = re.compile(r'--(?P<option_name>[a-zA-Z_]+)=(?P<value>[\S]+)')
    valid_options = [
        'dataset_path',
        'size_per_record',
        'output_path',
    ]
    settings = {
        # Default setting
        'dataset_path': 'images',
        'size_per_record': 52428800,    # 50MBytes
        'output_path': 'tfrecords'
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
                if option_name in ['size_per_record']:
                    option_value = int(option_value)
                settings[option_name] = option_value

    if settings['output_path'].endswith('/'):
        settings['output_path'] = settings['output_path'][:-1]
    
    return settings

# TFRecord 파일 생성 함수
def create_tfrecord(settings):
    if not os.path.isdir(settings['output_path']):
        os.mkdir(settings['output_path'], mode=0o777,)
    output_file_prefix = settings['output_path'] + '/out_'
    
    category_index = 0  # For one-hot encoding
    sample_counter = 0
    sample_size = 0     # Total sample size (Bytes)
    output_counter = 0
    
    writer = tf.io.TFRecordWriter(output_file_prefix + str(output_counter) + '.tfrecord')
    for category in os.listdir(settings['dataset_path']):
        category_dir = os.path.join(settings['dataset_path'], category)
        image_filenames = os.listdir(category_dir)
        
        for image_filename in image_filenames:
            image_path = os.path.join(category_dir, image_filename)
            image = tf.io.read_file(image_path)
            sample_counter += 1
            sample_size += os.path.getsize(image_path)

            example = tf.train.Example(features=tf.train.Features(feature={
                'name': tf.train.Feature(bytes_list=tf.train.BytesList(value=[image_filename.encode('utf-8')])),
                'category': tf.train.Feature(int64_list=tf.train.Int64List(value=[category_index])),
                'data': tf.train.Feature(bytes_list=tf.train.BytesList(value=[image.numpy()]))
            }))
            
            writer.write(example.SerializeToString())

            if sample_size >= settings['size_per_record']:
                writer.close()
                print('TFRecord file ' + output_file_prefix +  str(output_counter) + '.tfrecord (contains ' + str(sample_counter) + ' samples) created.')

                output_counter += 1
                sample_counter = 0
                sample_size = 0
                writer = tf.io.TFRecordWriter(output_file_prefix + str(output_counter) + '.tfrecord')

        category_index += 1
    
    if sample_counter != 0:
        writer.close()
        print('TFRecord file ' + output_file_prefix + str(output_counter) + '.tfrecord (contains ' + str(sample_counter) + ' samples) created.')

def main():
    settings = get_args(sys.argv[1:])
    create_tfrecord(settings)

main()