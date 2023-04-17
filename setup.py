from setuptools import setup
import sysconfig

setup(
    name='awscli-console-plugin',
    version='0.1',
    py_modules=['console'],
    data_files=[('/awscli-console-plugin', ['./awscli-console-plugin'])],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    url="https://github.com/b-b3rn4rd/awscli-console-plugin",
    author="Bernard Baltrusaitis",
    author_email="bernard@runawaylover.info",
    description="AWSCLI plugin to login to AWS Console using your IAM credentials",
    python_requires='>=3.6'
)