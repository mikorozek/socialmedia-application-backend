pipeline {
    agent any
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Test') {
            steps {
                script {
                    docker.image('python:3.9').inside {
                        sh 'pip install pytest'
                        sh 'pip install -r requirements.txt'
                        sh 'python -m pytest tests/'
                    }
                }
            }
        }
    }
}
