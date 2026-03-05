package com.longhe.learn.javaee.core.processor;

import com.google.auto.service.AutoService;
import com.longhe.learn.javaee.core.annotation.HelloWorld;
import com.squareup.javapoet.JavaFile;
import com.squareup.javapoet.MethodSpec;
import com.squareup.javapoet.TypeSpec;

import javax.annotation.processing.*;
import javax.lang.model.SourceVersion;
import javax.lang.model.element.Element;
import javax.lang.model.element.Modifier;
import javax.lang.model.element.TypeElement;
import javax.tools.Diagnostic;
import java.io.IOException;
import java.util.Set;

@AutoService(Processor.class)
@SupportedAnnotationTypes("com.longhe.learn.javaee.core.annotation.HelloWorld")
@SupportedSourceVersion(SourceVersion.RELEASE_17)
public class HelloWorldProcessor extends AbstractProcessor {

    @Override
    public boolean process(Set<? extends TypeElement> annotations, RoundEnvironment roundEnv) {
        for (Element element : roundEnv.getElementsAnnotatedWith(HelloWorld.class)) {
            String className = element.getSimpleName().toString();
            processingEnv.getMessager().printMessage(
                    Diagnostic.Kind.NOTE, "Processing @HelloWorld on: " + className
            );
            generateHelloWorldFactory(className);
        }
        return true;
    }

    private void generateHelloWorldFactory(String sourceClassName) {
        // Build the hello() method
        MethodSpec helloMethod = MethodSpec.methodBuilder("hello")
                .addModifiers(Modifier.PUBLIC, Modifier.STATIC)
                .returns(String.class)
                .addStatement("return $S", "Hello from " + sourceClassName + "!")
                .build();

        // Build the class
        TypeSpec factoryClass = TypeSpec.classBuilder("HelloWorldFactory")
                .addModifiers(Modifier.PUBLIC, Modifier.FINAL)
                .addMethod(helloMethod)
                .build();

        // Write the file
        JavaFile javaFile = JavaFile.builder("com.longhe.learn.generated", factoryClass)
                .build();

        try {
            javaFile.writeTo(processingEnv.getFiler());
        } catch (IOException e) {
            processingEnv.getMessager().printMessage(
                    Diagnostic.Kind.ERROR, "Failed to generate file: " + e.getMessage()
            );
        }
    }
}