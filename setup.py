from setuptools import setup

setup(
    name='awscli-console-plugin',
    version='0.1',
    py_modules=['console'],
    data_files=[('bin', ['./awscli-console-plugin'])],
    install_requires=[
        'awscli',
    ],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    url="https://github.com/b-b3rn4rd/awscli-console-plugin",
    author="Bernard Baltrusaitis",
    author_email="bernard@runawaylover.info",
    description="AWSCLI plugin to login to AWS Console using your IAM credentials",
    python_requires='>=3.6',
)