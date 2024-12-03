pipeline {
    agent any

    environment {
        NEXUS_URL = 'http://<IP_MASZYNY>:8081'  // Adres Twojego Nexusa
        NEXUS_REPO = 'python-repo'  // Repozytorium w Nexusie, gdzie będziesz wysyłać artefakty
        NEXUS_CREDENTIALS = 'nexus-admin'  // ID poświadczeń w Jenkinsie
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Install dependencies') {
            steps {
                sh 'pip install -r requirements.txt'
            }
        }

        stage('Build') {
            steps {
                sh 'python setup.py sdist bdist_wheel'
            }
        }

        stage('Upload to Nexus') {
            steps {
                script {
                    def files = findFiles(glob: 'dist/*')
                    files.each {
                        sh "twine upload --repository-url http://<IP_MASZYNY>:8081/repository/python-repo/ ${it}"
                    }
                }
            }
        }
    }
}
