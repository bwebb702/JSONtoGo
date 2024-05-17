import tkinter as tk
import sys
from tkinter import filedialog

root = tk.Tk()
root.withdraw()

file_path = filedialog.askopenfilename()
sys.stdout.write(file_path)