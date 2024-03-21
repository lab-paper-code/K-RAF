import sys
import re
import csv

strace_log_path = sys.argv[1]
line_parser = re.compile(
    r'\A(?P<process_id>[0-9]+)\s+(?P<time>[0-9\.]+)\s+(?P<syscall_name>[^(]+)\((?P<arguments>[\s\S]+)\)\s*=\s*(?P<return_value>[0-9\-]+)[\s\S]*\Z')
argument_parser = re.compile(r'([^,]+),? ?')
current_stats = {}  # fd: offset
fd_file_map = {}  # fd: file path
timeline = {}  # Time (누적): (file path, syscall, offset)
time = 0

with open(strace_log_path, 'r') as f:
    log_contents = f.readlines()
    for line in log_contents:
        parsing_result = line_parser.search(line)
        if parsing_result is None:
            continue
        time += float(parsing_result.group('time'))
        syscall = parsing_result.group('syscall_name')
        args = parsing_result.group('arguments')
        return_value = parsing_result.group('return_value')
        if syscall == 'openat':
            fd = return_value
            file_path = argument_parser.findall(args)[1].replace('"', '')
            current_stats[fd] = 0
            fd_file_map[fd] = file_path
        elif syscall == 'read':
            fd = argument_parser.findall(args)[0]
            if fd not in current_stats.keys():
                continue
            read_bytes = int(return_value)
            current_stats[fd] += read_bytes
        elif syscall == 'lseek':
            fd = argument_parser.findall(args)[0]
            if fd not in current_stats.keys():
                continue
            after_offset = int(return_value)
            current_stats[fd] = after_offset
        else:
            continue
        timeline[time] = (fd_file_map[fd], syscall, current_stats[fd])

with open('./opened_files/' + strace_log_path.split('.')[0] + '_timeline.csv', 'w+') as f:
    csv_writer = csv.writer(f, delimiter=',')
    csv_writer.writerow(['Time', 'File Path', 'Syscall', 'Offset'])
    for time, read_info in timeline.items():
        file_path, syscall, offset = read_info
        csv_writer.writerow([time, file_path, syscall, offset])
