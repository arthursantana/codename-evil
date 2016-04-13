var gulp = require('gulp');
var watch = require('gulp-watch');
//var uglify = require('gulp-uglify');

var sass = require('gulp-sass');
var react = require('gulp-react');


gulp.task('html', function () {
   gulp.src(['client/*.htm'])
   .pipe(gulp.dest('dist/static'));
});

gulp.task('css', function () {
   gulp.src('client/scss/*.scss')
   .pipe(sass())
   .pipe(gulp.dest('dist/static/css'))
});

gulp.task('img', function () {
   gulp.src(['art/buildings/*.png', 'art/resources/*.png'])
   .pipe(gulp.dest('dist/static/img'))
});

gulp.task('js', function () {
   gulp.src('client/js/*.js')
   .pipe(react())
   .pipe(gulp.dest('dist/static/js'));
});

gulp.task('libs', function () {
   gulp.src('node_modules/zepto-full/zepto.js')
   //gulp.src('node_modules/zepto-full/zepto.min.js')
   .pipe(gulp.dest('dist/static/js'));

   gulp.src('node_modules/react/dist/react.js')
   //gulp.src('node_modules/react/dist/react.min.js')
   .pipe(gulp.dest('dist/static/js'));

   gulp.src('node_modules/react-dom/dist/react-dom.js')
   //gulp.src('node_modules/react-dom/dist/react-dom.min.js')
   .pipe(gulp.dest('dist/static/js'));
});

gulp.task('default', ['html', 'css', 'img', 'js', 'libs']);

gulp.task('watch', function () {
   gulp.watch(['client/*', 'client/*/*', 'gulpfile.js'], ['default'])
});

//gulp.task('node_modules', function () {
//   gulp.src('node_modules/react/dist/react.js')
//   .pipe(gulp.dest('dist/static/js'));
//
//   gulp.src(['node_modules/bootstrap/dist/css/bootstrap.min.css', 'node_modules/font-awesome/css/font-awesome.min.css'])
//   .pipe(gulp.dest('dist/static/css'));
//});

