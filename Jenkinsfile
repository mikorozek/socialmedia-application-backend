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
                sh '''
                    python3 -m pip install --user pytest
                    python3 -m pip install --user -r requirements.txt
                    python3 -m pytest tests/
                '''
            }
        }
    }
}
