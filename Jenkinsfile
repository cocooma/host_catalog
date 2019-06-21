#!/usr/bin/env groovy
pipeline {
  agent {
    label 'default-slave'
  }

  stages {
    stage('Build') {
      steps {
        ansiColor('xterm') {
          sh("make")
        }
      }
    }
  }
}
