from setuptools import setup, find_packages

setup(
    name='example-python-app',  # Nazwa aplikacji
    version='1.0',  # Wersja aplikacji
    packages=find_packages(),  # Automatycznie wykrywa wszystkie pakiety w projekcie
    install_requires=[  # Lista zależności, które muszą być zainstalowane
        'requests',  # Przykładowa zależność
    ],
    author='Your Name',  # Twoje imię
    author_email='your.email@example.com',  # Twój e-mail
    description='An example Python app',  # Krótkie streszczenie aplikacji
    long_description=open('README.md').read(),  # Długie opis aplikacji
    long_description_content_type='text/markdown',  # Typ treści README
    url='http://example.com',  # URL do strony aplikacji (jeśli jest)
    classifiers=[  # Klasyfikatory, które pomagają w wyszukiwaniu aplikacji
        'Programming Language :: Python :: 3',
        'License :: OSI Approved :: MIT License',
    ],
)
