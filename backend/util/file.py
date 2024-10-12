import os
import re
from typing import List, Tuple

def read_lines_from_directory(directory):
    try:
        files = os.listdir(directory)
    except OSError as e:
        print(f"Error reading directory {directory}: {e}")
        return None, e

    # Regular expression to extract numbers from filenames
    regex = re.compile(r'\d+')

    # List to store filenames and their corresponding numbers
    file_list = []

    for file in files:
        if file.endswith('.txt'):
            matches = regex.findall(file)
            if matches:
                try:
                    number = int(matches[0])
                    file_list.append((file, number))
                except ValueError:
                    continue

    # Sort files based on the extracted number
    file_list.sort(key=lambda x: x[1])

    lines = []
    for file, _ in file_list:
        try:
            with open(os.path.join(directory, file), 'r') as f:
                lines.append(f.read())
        except OSError as e:
            print(f"Error reading file {file}: {e}")
            continue

    return lines, None

def read_files_from_directory(dir_path: str) -> List[os.DirEntry]:
    """
    Reads all files from the specified directory, sorts them by numeric order
    extracted from their filenames, and returns them as a list of os.DirEntry objects.
    
    :param dir_path: Path to the directory containing the files.
    :return: List of os.DirEntry objects sorted by numeric order.
    """
    # Regular expression to extract numbers from filenames
    number_re = re.compile(r'\d+')
    
    # List to store tuples of (os.DirEntry, number)
    file_list = []
    
    # Iterate over directory entries
    with os.scandir(dir_path) as entries:
        for entry in entries:
            if entry.is_file():
                # Extract number from filename
                match = number_re.search(entry.name)
                if match:
                    number = int(match.group())
                    file_list.append((entry, number))
    
    # Sort files by the extracted number
    file_list.sort(key=lambda x: x[1])
    
    # Extract sorted os.DirEntry objects
    sorted_files = [entry.name for entry, _ in file_list]
    
    return sorted_files

def save_list_to_files(input_list, path, offset):
    try:
        for i, line in enumerate(input_list):
            file_path = f"{path}{i + offset}.txt"
            with open(file_path, 'w') as file:
                file.write(line)
    except Exception as e:
        return e
    return None
