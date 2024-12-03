pipeline {
    agent any

    environment {
        NEXUS_URL = 'http://<IP_MASZYNY>:8081'  // Adres Twojego Nexusa
        NEXUS_REPO = 'example-maven'  // Repozytorium w Nexusie, gdzie będziesz wysyłać artefakty
        NEXUS_CREDENTIALS = 'nexus-admin'  // ID poświadczeń w Jenkinsie
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build') {
            steps {
                sh 'mvn clean package'
            }
        }

        stage('Upload to Nexus') {
            steps {
                nexusArtifactUploader artifacts: [[artifactId: 'example-app',
                                                   classifier: '',
                                                   file: 'target/example-app-1.0.jar',
                                                   type: 'jar']],
                                        credentialsId: "${env.NEXUS_CREDENTIALS}",
                                        groupId: 'com.example',
                                        nexusUrl: "${env.NEXUS_URL}",
                                        nexusVersion: 'nexus3',
                                        protocol: 'http',
                                        repository: "${env.NEXUS_REPO}",
                                        version: '1.0'
            }
        }
    }
}
