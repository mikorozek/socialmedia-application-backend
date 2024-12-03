pipeline {
    agent any
    
    environment {
        NEXUS_URL = 'http://<IP_MASZYNY>:8081'  // Adres Twojego Nexusa
        NEXUS_REPO = 'python-repo'  // Repozytorium w Nexusie, gdzie będziesz wysyłać artefakty
        NEXUS_CREDENTIALS = 'nexus-admin'  // ID poświadczeń w Jenkinsie
    }

    stages {
        // Etap pobierania kodu z repozytorium
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        // Etap testowania aplikacji w kontenerze Docker z Pythonem
        stage('Test') {
            agent {
                docker {
                    image 'python:3.9'  // Obraz Docker z Pythonem 3.9
                }
            }
            steps {
                withEnv(["HOME=${env.WORKSPACE}"]) {
                    // Instalacja zależności z requirements.txt
                    sh 'pip install -r requirements.txt'
                    
                    // Instalacja pytest (jeśli nie ma w requirements.txt)
                    sh 'pip install pytest'
                    
                    // Uruchomienie testów
                    sh 'python -m pytest tests/'
                }
            }
        }

        // Etap budowania artefaktów Pythona
        stage('Build') {
            agent {
                docker {
                    image 'python:3.9'  // Obraz Docker z Pythonem 3.9
                }
            }
            steps {
                withEnv(["HOME=${env.WORKSPACE}"]) {
                    // Budowanie pakietu Python w formacie sdist i bdist_wheel
                    sh 'python setup.py sdist bdist_wheel'
                }
            }
        }

        // Etap przesyłania artefaktów do Nexusa
        stage('Upload to Nexus') {
            agent {
                docker {
                    image 'python:3.9'  // Obraz Docker z Pythonem 3.9
                }
            }
            steps {
                script {
                    // Znalezienie plików w folderze dist
                    def files = findFiles(glob: 'dist/*')
                    files.each { file ->
                        // Wysyłanie każdego pliku do repozytorium Nexus
                        sh "twine upload --repository-url ${NEXUS_URL}/repository/${NEXUS_REPO}/ ${file}"
                    }
                }
            }
        }
    }
}
